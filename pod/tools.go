package pod

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/transport/spdy"
)

type Container struct {
	Name  string
	Image string
}
type PodController struct {
	metav1.OwnerReference
}

// IsReady check whether the pod is ready.
func (h *Handler) IsReady(name string) bool {
	pod, err := h.Get(name)
	if err != nil {
		return false
	}
	// 必须要 type=Ready 和 Status=True 才能算 Pod 已经就绪
	for _, cond := range pod.Status.Conditions {
		if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

// WaitReady waiting for the pod to be in the ready status.
func (h *Handler) WaitReady(name string) error {
	if h.IsReady(name) {
		return nil
	}

	errCh := make(chan error, 1)
	chkCh := make(chan struct{}, 1)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	ctxCheck, cancelCheck := context.WithCancel(h.ctx)
	ctxWatch, cancelWatch := context.WithCancel(h.ctx)
	defer cancelCheck()
	defer cancelWatch()

	// this goroutine used to check whether pod is ready and exists.
	// if pod already ready, return nil.
	// if pod does not exist, return error.
	// if pod exists but not ready, return nothing.
	go func(context.Context) {
		<-chkCh
		for i := 0; i < 6; i++ {
			// pod is already ready, return nil
			if h.IsReady(name) {
				errCh <- nil
				return
			}
			// pod no longer exists, return err
			_, err := h.Get(name)
			if k8serrors.IsNotFound(err) {
				errCh <- err
				return
			}
			// pod exists but not ready, return northing.
			if err == nil {
				return
			}
			log.Error(err)
			time.Sleep(time.Second * 10)
		}
	}(ctxCheck)

	// this goroutine used to watch pod.
	go func(ctx context.Context) {
		for {
			timeout := int64(0)
			listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
			listOptions.TimeoutSeconds = &timeout
			watcher, err := h.clientset.CoreV1().Pods(h.namespace).Watch(h.ctx, listOptions)
			if err != nil {
				errCh <- err
				return
			}
			chkCh <- struct{}{}
			for event := range watcher.ResultChan() {
				switch event.Type {
				case watch.Modified:
					if h.IsReady(name) {
						watcher.Stop()
						errCh <- nil
						return
					}
				case watch.Deleted:
					watcher.Stop()
					errCh <- fmt.Errorf("pod/%s was deleted", name)
					return
				case watch.Bookmark:
					log.Debug("watch pod: bookmark")
				case watch.Error:
					log.Debug("watch pod: error")
				}
			}
			// If event channel is closed, it means the kube-apiserver has closed the connection.
			log.Debug("watch pod: reconnect to kubernetes")
			watcher.Stop()
		}
	}(ctxWatch)

	select {
	case sig := <-sigCh:
		return fmt.Errorf("canceled by signal: %v", sig.String())
	case err := <-errCh:
		return err
	}
}

//// WaitReady waiting for the pod to be in the ready status.
//func (h *Handler) WaitReady2(name string) (err error) {
//    var (
//        watcher watch.Interface
//        timeout = int64(0)
//    )
//    // 在 watch 之前就先判断 pod 是否就绪, 如果就绪了就没必要 watch 了
//    if h.IsReady(name) {
//        return nil
//    }
//    // 判断 pod 是否存在
//    if _, err := h.Get(name); err != nil {
//        return err
//    }
//    for {
//        // pod 没有就绪, 那么就开始监听 pod 的事件
//        listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
//        listOptions.TimeoutSeconds = &timeout
//        watcher, err = h.clientset.CoreV1().Pods(h.namespace).Watch(h.ctx, listOptions)
//        if err != nil {
//            return err
//        }
//        for event := range watcher.ResultChan() {
//            switch event.Type {
//            case watch.Modified:
//                if h.IsReady(name) {
//                    watcher.Stop()
//                    return nil
//                }
//            case watch.Deleted:
//                watcher.Stop()
//                return fmt.Errorf("%s deleted", name)
//            case watch.Bookmark:
//                log.Debug("watch pod: bookmark")
//            case watch.Error:
//                log.Debug("watch pod: error")
//            }
//        }
//        // watcher 因为 keepalive 超时断开了连接, 关闭了 channel
//        log.Debug("watch pod: reconnect to kubernetes")
//        watcher.Stop()
//    }
//}

// GetUID returns pod uuid
func (h *Handler) GetUID(object interface{}) (string, error) {
	switch val := object.(type) {
	// if object type is string, the object is regarded as pod name,
	// and check whether the pod exists.
	case string:
		pod, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return string(pod.ObjectMeta.UID), nil
	case *corev1.Pod:
		return string(val.ObjectMeta.UID), nil
	case corev1.Pod:
		return string(val.ObjectMeta.UID), nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetIP returns pod ip.
func (h *Handler) GetIP(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pod, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return pod.Status.PodIP, nil
	case *corev1.Pod:
		return val.Status.PodIP, nil
	case corev1.Pod:
		return val.Status.PodIP, nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetNodeIP get the ip addr of the node where pod is located.
func (h *Handler) GetNodeIP(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pod, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return pod.Status.HostIP, nil
	case *corev1.Pod:
		return val.Status.HostIP, nil
	case corev1.Pod:
		return val.Status.HostIP, nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetNodeName get the name of the node where pod is located.
func (h *Handler) GetNodeName(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pod, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return pod.Spec.NodeName, nil
	case *corev1.Pod:
		return val.Spec.NodeName, nil
	case corev1.Pod:
		return val.Spec.NodeName, nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetAge get the pod age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		pod, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		ctime := pod.ObjectMeta.CreationTimestamp.Time
		return time.Now().Sub(ctime), nil
	case *corev1.Pod:
		ctime := val.ObjectMeta.CreationTimestamp.Time
		return time.Now().Sub(ctime), nil
	case corev1.Pod:
		ctime := val.ObjectMeta.CreationTimestamp.Time
		return time.Now().Sub(ctime), nil
	default:
		return time.Duration(int64(0)), ErrInvalidToolsType
	}

}

// GetStatus get the status of the pod.
// There are the valid statuses of the pod
// Pending:     pod has been accepted by the system, but one or more of the
//              containers has not been started.
// Running:     pod is running and all of the containers have been started.
// Succeeded:   all containers in the pod have voluntarily terminated.
// Failed:      all containers in the pod have terminated, and at least one
//              container hasterminated in a failure (exited with a non-zero
//              exit code or was stopped by the system).
// Unknown:     for some reason the state of the pod could not be obtained,
//              typically due to an error in communicating with the host of the pod.
func (h *Handler) GetStatus(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pod, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return string(pod.Status.Phase), nil
	case *corev1.Pod:
		return string(val.Status.Phase), nil
	case corev1.Pod:
		return string(val.Status.Phase), nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetQosClass get the "Quality of Service" of the pod.
// PodQOSGuaranteed is the Guaranteed qos class.
// PodQOSBurstable is the Burstable qos class.
// PodQOSBestEffort is the BestEffort qos class.
func (h *Handler) GetQosClass(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		pod, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return string(pod.Status.QOSClass), nil
	case *corev1.Pod:
		return string(val.Status.QOSClass), nil
	case corev1.Pod:
		return string(val.Status.QOSClass), nil
	default:
		return "", ErrInvalidToolsType
	}
}

// GetPVC get all persistentvolumeclaim mounted by this pod.
func (h *Handler) GetPVC(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		pod, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getPVC(pod), nil
	case *corev1.Pod:
		return h.getPVC(val), nil
	case corev1.Pod:
		return h.getPVC(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getPVC(pod *corev1.Pod) []string {
	var pl []string
	for _, volume := range pod.Spec.Volumes {
		// 要先判断 volume.PersistentVolumeClaim 是否为空, 如果不判断而直接获取
		// volume.PersistentVolumeClaim.ClaimName 相当于操纵值为 nil 的指针(空指针),
		// 程序会直接 panic.
		if volume.PersistentVolumeClaim != nil {
			pl = append(pl, volume.PersistentVolumeClaim.ClaimName)
		}
	}
	return pl
}

// GetPV get all persistentvolume mounted by this pod.
func (h *Handler) GetPV(object interface{}) ([]string, error) {
	// It does not need to check whether the pod exists,
	// GetPVC will do check.
	pvcList, err := h.GetPVC(object)
	if err != nil {
		return nil, err
	}

	var pl []string
	for _, pvc := range pvcList {
		pvcObj, err := h.clientset.CoreV1().
			PersistentVolumeClaims(h.namespace).Get(h.ctx, pvc, h.Options.GetOptions)
		if err == nil {
			pl = append(pl, pvcObj.Spec.VolumeName)
		}
	}
	return pl, nil
}

// GetController get the owner of the pod.
// Controller maybe deployment, statefulset, daemonset, job, replicaset, rc.
func (h *Handler) GetController(object interface{}) (*PodController, error) {
	switch val := object.(type) {
	case string:
		pod, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getController(pod)
	case *corev1.Pod:
		return h.getController(val)
	case corev1.Pod:
		return h.getController(&val)
	default:
		return nil, ErrInvalidToolsType
	}

}
func (h *Handler) getController(pod *corev1.Pod) (*PodController, error) {
	pc := new(PodController)
	ownerRef := metav1.GetControllerOf(pod)
	if ownerRef == nil {
		return pc, fmt.Errorf("the pod %q doesn't have controller", pod.Name)
	}
	return &PodController{OwnerReference: *ownerRef}, nil
}

// GetContainers get all containers of the pod.
func (h *Handler) GetContainers(object interface{}) ([]Container, error) {

	switch val := object.(type) {
	case string:
		pod, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getContainers(pod), nil
	case *corev1.Pod:
		return h.getContainers(val), nil
	case corev1.Pod:
		return h.getContainers(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getContainers(pod *corev1.Pod) []Container {
	var cl []Container
	// pod.Status.ContainerStatuses 这个是就绪之后才会有的,不要从这里获取 containers
	//for _, cs := range pod.Status.ContainerStatuses {
	//    c := Container{
	//        Name:  cs.Name,
	//        Image: cs.Image,
	//    }
	//    cl = append(cl, c)
	//}
	for _, container := range pod.Spec.Containers {
		c := Container{
			Name:  container.Name,
			Image: container.Image,
		}
		cl = append(cl, c)
	}
	return cl
}

// GetInitContainers get all init containers of the pod.
func (h *Handler) GetInitContainers(object interface{}) ([]Container, error) {
	switch val := object.(type) {
	case string:
		pod, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getInitContainers(pod), nil
	case *corev1.Pod:
		return h.getInitContainers(val), nil
	case corev1.Pod:
		return h.getInitContainers(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getInitContainers(pod *corev1.Pod) []Container {
	var cl []Container
	for _, container := range pod.Spec.InitContainers {
		c := Container{
			Name:  container.Name,
			Image: container.Image,
		}
		cl = append(cl, c)
	}
	return cl
}

// GetReadyContainers get all ready containers of the pod.
func (h *Handler) GetReadyContainers(object interface{}) ([]Container, error) {
	switch val := object.(type) {
	case string:
		pod, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getReadyContainers(pod), nil
	case *corev1.Pod:
		return h.getReadyContainers(val), nil
	case corev1.Pod:
		return h.getReadyContainers(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getReadyContainers(pod *corev1.Pod) []Container {
	var cl []Container
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.Ready {
			c := Container{
				Name:  cs.Name,
				Image: cs.Image,
			}
			cl = append(cl, c)
		}
	}
	return cl
}

// references:
//    https://miminar.fedorapeople.org/_preview/openshift-enterprise/registry-redeploy/go_client/executing_remote_processes.html
//    https://stackoverflow.com/questions/43314689/example-of-exec-in-k8ss-pod-by-using-go-client
//    https://github.com/kubernetes/kubernetes/blob/v1.6.1/test/e2e/framework/exec_util.go
//    https://github.com/kubernetes/client-go/issues/464  (How to make a web terminal)
//    https://github.com/kubernetes/dashboard/blob/master/src/app/backend/handler/terminal.go
//    http://maoqide.live/post/cloud/kubernetes-webshell/

// client-go 提供的 k8s.io/client-go/tools/remotecommand 包, 提供了方法与集群中的
// 容器进行长连接,并设置容器的 stdin, stdout 等.
// remotecommand 包提供了基于 SPDY 协议的 Executor interface 进行和 pod 终端的流
// 的传输. 初始化一个 Executor 很简单,只需要调用 NewSPDYExecutor 并传入对应参数.
// Executor 的 Stream 方法, 会建立一个传输的连接,直到服务端和调用端一段关闭连接,
// 才会停止传输. 常用的做法是定义一个 PtyHandler 的 interface, 然后使用你想用的
// 客户端实现该 interface 对应的 Read(p []byte) (int, error) 和 Write(p []byte) (int, error)
// 方法即可, Executor 调用 Stream 方法时,只要将 StreamOptions 的 Stdin Stdout
// 都设置为 PtyHandler, Executor 就会通过你定义的 write 和 read 方法来传输数据.

// Execute will executing remote processes in a container of the pod.
// If no container name is specified, Execute will executing a process
// in the first container of the pod by default.
// It will returns error, If the pod not ready. It's your responsibility to ensure
// the pod Is running and ready.
//
// The remote processes default stdin, stdou, stderr are os.Stdin, os.Stdout, os.Stderr.
func (h *Handler) Execute(podName, containerName string, command []string) error {
	// if pod not found, returns error.
	pod, err := h.Get(podName)
	if err != nil {
		return err
	}

	// if containerName is empty, execute command in first container of the pod.
	if len(containerName) == 0 {
		containerName = pod.Spec.Containers[0].Name
	}

	// Prepare the API URL used to execute another process within the Pod.  In
	// this case, we'll run a remote shell.
	req := h.restClient.Post().
		Namespace(h.namespace).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(h.config, "POST", req.URL())
	if err != nil {
		return err
	}

	//// Put the terminal into raw mode to prevent it echoing characters twice.
	//// The integer file descriptor associated with the stream stdin, stdout
	//// and stderr are 0, 1 and 2, respectively.
	////oldState, err := terminal.MakeRaw(0)
	////defer terminal.Restore(0, oldState)
	//oldState, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	//if err != nil {
	//    return fmt.Errorf("Failed to set raw mod on Stdin: %v\n", err)
	//}
	//defer terminal.Restore(int(os.Stdin.Fd()), oldState)

	// if passed ptyhandler is nil
	return exec.Stream(remotecommand.StreamOptions{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    true,
	})
}

// ExecuteWithPty will executing remote processes in a container of the pod.
// If no container name is specified, Execute will executing a process
// in the first container of the pod by default.
// It will returns error, If the pod not ready. It's your responsibility to ensure
// the pod Is running and ready.
//
// You should provide a PtyHandler interface.
// What is pty, please refer to https://man7.org/linux/man-pages/man7/pty.7.html
func (h *Handler) ExecuteWithPty(podName, containerName string, command []string, pty PtyHandler) error {
	// if pod not found, returns error.
	pod, err := h.Get(podName)
	if err != nil {
		return err
	}

	// if containerName is empty, execute command in first container of the pod.
	if len(containerName) == 0 {
		containerName = pod.Spec.Containers[0].Name
	}

	// Prepare the API URL used to execute another process within the Pod.
	// In this case, we'll run a remote shell.
	req := h.restClient.Post().
		Namespace(h.namespace).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(h.config, "POST", req.URL())
	if err != nil {
		return err
	}

	// Connect the process std(in,out,err) to the remote shell process.
	if pty == nil || reflect.ValueOf(pty).IsNil() {
		return executor.Stream(remotecommand.StreamOptions{
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Tty:    true,
		})
	}
	return executor.Stream(remotecommand.StreamOptions{
		Stdin:             pty,
		Stdout:            pty,
		Stderr:            pty,
		TerminalSizeQueue: pty,
		Tty:               true,
	})
}

// ExecuteWithStream will executing remote processes in a container of the pod.
// If no container name is specified, Execute will executing a process
// in the first container of the pod by default.
// It will returns error, If the pod not ready. It's your responsibility to ensure
// the pod Is running and ready.
//
// You should manually specify that the stdin, stdout and stderr of the remote shell process.
func (h *Handler) ExecuteWithStream(podName, containerName string, command []string, stdin io.Reader, stdout, stderr io.Writer) error {
	// if pod not found, returns error.
	pod, err := h.Get(podName)
	if err != nil {
		return err
	}

	// if containerName is empty, execute command in first container of the pod.
	if len(containerName) == 0 {
		containerName = pod.Spec.Containers[0].Name
	}

	// Prepare the API URL used to execute another process within the Pod.
	// In this case, we'll run a remote shell.
	req := h.restClient.Post().
		Namespace(h.namespace).
		Resource("pods").
		Name(podName).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   command,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(h.config, "POST", req.URL())
	if err != nil {
		return err
	}

	// Connect the process std(in,out,err) to the remote shell process.
	return executor.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    true,
	})
}

// PortForward forward a local port to the pod.
//
// ref:
//   https://github.com/kubernetes/client-go/issues/51
//   https://github.com/anthhub/forwarder
//   https://www.modb.pro/db/137716
func (h *Handler) PortForward(podName string, addr net.IP, localPort, remotePort int32) error {
	roundTripper, upgrader, err := spdy.RoundTripperFor(h.config)
	if err != nil {
		return err
	}

	_, _ = roundTripper, upgrader
	return nil
}
