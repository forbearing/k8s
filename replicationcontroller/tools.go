package replicationcontroller

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// IsReady check if the replicationcontroller is ready
func (h *Handler) IsReady(name string) bool {
	checkGeneration := func(rc *corev1.ReplicationController) bool {
		if rc.Generation != rc.Status.ObservedGeneration {
			return false
		}
		return true
	}
	checkReplicas := func(rc *corev1.ReplicationController) bool {
		if rc.Spec.Replicas == nil {
			return false
		}
		if *rc.Spec.Replicas != rc.Status.ReadyReplicas {
			return false
		}
		if *rc.Spec.Replicas != rc.Status.AvailableReplicas {
			return false
		}
		if *rc.Spec.Replicas != rc.Status.Replicas {
			return false
		}
		if *rc.Spec.Replicas != rc.Status.FullyLabeledReplicas {
			return false
		}
		return true
	}

	rc, err := h.Get(name)
	if err != nil {
		return false
	}
	if checkGeneration(rc) && checkReplicas(rc) {
		return true
	}
	return false
}

// WaitReady waiting for the replicationcontroller to be in the ready status.
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

	// this goroutine used to check whether replicationcontroller is ready and exists.
	// if replicationcontroller already ready, return nil.
	// if replicationcontroller does not exist, return error.
	// if replicationcontroller exists but not ready, return nothing.
	go func(context.Context) {
		<-chkCh
		for i := 0; i < 6; i++ {
			// replicationcontroller is already ready, return nil
			if h.IsReady(name) {
				errCh <- nil
				return
			}
			// replicationcontroller no longer exists, return err
			_, err := h.Get(name)
			if k8serrors.IsNotFound(err) {
				errCh <- err
				return
			}
			// replicationcontroller exists but not ready, return nothing.
			if err == nil {
				return
			}
			log.Error(err)
			time.Sleep(time.Second * 10)
		}
	}(ctxCheck)

	// this goroutine used to watch replicationcontroller.
	go func(ctx context.Context) {
		for {
			timeout := int64(0)
			listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
			listOptions.TimeoutSeconds = &timeout
			watcher, err := h.clientset.CoreV1().ReplicationControllers(h.namespace).Watch(h.ctx, listOptions)
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
					errCh <- fmt.Errorf("replicationcontroller/%s was deleted", name)
					return
				case watch.Bookmark:
					log.Debug("watch replicationcontroller: bookmark")
				case watch.Error:
					log.Debug("watch replicationcontroller: error")
				}
			}
			// If event channel is closed, it means the kube-apiserver has closed the connection.
			log.Debug("watch replicationcontroller: reconnect to kubernetes")
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

//// WaitReady wait for the replicationcontroller to be in the ready status
//func (h *Handler) WaitReady(name string, check bool) (err error) {
//    var (
//        watcher watch.Interface
//        timeout = int64(0)
//    )
//    // 在 watch 之前就先判断 replicationcontroller 是否就绪, 如果就绪了就没必要 watch 了
//    if h.IsReady(name) {
//        return
//    }
//    // 是否判断 replicationcontroller 是否存在
//    if check {
//        if _, err = h.Get(name); err != nil {
//            return
//        }
//    }
//    for {
//        // replicationcontroller 没有就绪, 那么就开始监听 replicationcontroller 的事件
//        listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
//        listOptions.TimeoutSeconds = &timeout
//        watcher, err = h.clientset.CoreV1().ReplicationControllers(h.namespace).Watch(h.ctx, listOptions)
//        for event := range watcher.ResultChan() {
//            switch event.Type {
//            case watch.Modified:
//                if h.IsReady(name) {
//                    watcher.Stop()
//                    return
//                }
//            case watch.Deleted:
//                watcher.Stop()
//                return fmt.Errorf("%s deleted", name)
//            case watch.Bookmark:
//                log.Debug("watch replicationcontroller: bookmark")
//            case watch.Error:
//                log.Debug("watch replicationcontroller: error")
//            }
//        }
//        // watcher 因为 keepalive 超时断开了连接, 关闭了 channel
//        log.Debug("watch replicationcontroller: reconnect to kubernetes")
//        watcher.Stop()
//    }
//}

// GetPods
// GetPVC
// GetPV
// GetAge

// GetPods get replicationcontroller all pods
func (h *Handler) GetPods(object interface{}) ([]*corev1.Pod, error) {
	switch val := object.(type) {
	case string:
		rc, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getPods(rc)
	case *corev1.ReplicationController:
		return h.getPods(val)
	case corev1.ReplicationController:
		return h.getPods(&val)
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getPods(rc *corev1.ReplicationController) ([]*corev1.Pod, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	podList, err := h.clientset.CoreV1().Pods(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}

	var pl []*corev1.Pod
	for i := range podList.Items {
		for _, or := range podList.Items[i].OwnerReferences {
			if or.Name == rc.Name {
				pl = append(pl, &podList.Items[i])
			}
		}
	}
	return pl, nil
}

// GetPVC get all persistentvolumeclaims mounted by the replicationcontroller.
func (h *Handler) GetPVC(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		rc, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getPVC(rc), nil
	case *corev1.ReplicationController:
		return h.getPVC(val), nil
	case corev1.ReplicationController:
		return h.getPVC(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getPVC(rc *corev1.ReplicationController) []string {
	var pl []string
	for _, volume := range rc.Spec.Template.Spec.Volumes {
		if volume.PersistentVolumeClaim != nil {
			pl = append(pl, volume.PersistentVolumeClaim.ClaimName)
		}
	}
	return pl
}

// GetPV get all persistentvolumes mounted by the replicationcontroller.
func (h *Handler) GetPV(object interface{}) ([]string, error) {
	// GetPV does not need to check whether replicationcontroller is exists.
	// GetPVC will do it.
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

// GetAge returns replicationcontroller age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		rc, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		return time.Now().Sub(rc.CreationTimestamp.Time), nil
	case *corev1.ReplicationController:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case corev1.ReplicationController:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}
