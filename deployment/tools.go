package deployment

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// IsReady check if the deployment is ready.
// ref: https://github.com/kubernetes/kubernetes/blob/a1128e380c2cf1c2d7443694673d9f1dd63eb518/staging/src/k8s.io/kubectl/pkg/polymorphichelpers/rollout_status.go#L59
func (h *Handler) IsReady(name string) bool {
	checkGeneration := func(deploy *appsv1.Deployment) bool {
		if deploy.Generation != deploy.Status.ObservedGeneration {
			return false
		}
		return true
	}
	checkReplicas := func(deploy *appsv1.Deployment) bool {
		if deploy.Spec.Replicas == nil {
			return false
		}
		if *deploy.Spec.Replicas != deploy.Status.ReadyReplicas {
			return false
		}
		if *deploy.Spec.Replicas != deploy.Status.AvailableReplicas {
			return false
		}
		if *deploy.Spec.Replicas != deploy.Status.UpdatedReplicas {
			return false
		}
		if *deploy.Spec.Replicas != deploy.Status.Replicas {
			return false
		}
		return true
	}
	checkCondition := func(deploy *appsv1.Deployment) bool {
		for _, cond := range deploy.Status.Conditions {
			if cond.Type == appsv1.DeploymentAvailable && cond.Status == corev1.ConditionTrue {
				return true
			}
		}
		return false
	}

	deploy, err := h.Get(name)
	if err != nil {
		return false
	}
	if checkGeneration(deploy) && checkReplicas(deploy) && checkCondition(deploy) {
		return true
	}
	return false
}

// WaitReady waiting for the deployment to be in the ready status.
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

	// this goroutine used to check whether deployment is ready and exists.
	// if deployment already ready, return nil.
	// if deployment does not exist, return error.
	// if deployment exists but not ready, return nothing.
	go func(context.Context) {
		<-chkCh
		for i := 0; i < 6; i++ {
			// deployment is already ready, return nil
			if h.IsReady(name) {
				errCh <- nil
				return
			}
			// deployment no longer exists, return err
			_, err := h.Get(name)
			if k8serrors.IsNotFound(err) {
				errCh <- err
				return
			}
			// deployment exists but not ready, return nothing.
			if err == nil {
				return
			}
			log.Error(err)
			time.Sleep(time.Second * 10)
		}
	}(ctxCheck)

	// this goroutine used to watch deployment.
	go func(ctx context.Context) {
		for {
			timeout := int64(0)
			listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
			listOptions.TimeoutSeconds = &timeout
			watcher, err := h.clientset.AppsV1().Deployments(h.namespace).Watch(h.ctx, listOptions)
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
					errCh <- fmt.Errorf("deployment/%s was deleted", name)
					return
				case watch.Bookmark:
					log.Debug("watch deployment: bookmark")
				case watch.Error:
					log.Debug("watch deployment: error")
				}
			}
			// If event channel is closed, it means the kube-apiserver has closed the connection.
			log.Debug("watch deployment: reconnect to kubernetes")
			watcher.Stop()
		}
	}(ctxWatch)

	select {
	case sig := <-sigCh:
		return fmt.Errorf("cancelled by signal: %s", sig.String())
	case err := <-errCh:
		return err
	}
}

//// WaitReady waiting for the deployment to be in the ready state.
//func (h *Handler) WaitReady2(name string) (err error) {
//    var (
//        watcher watch.Interface
//        timeout = int64(0)
//    )
//    // 在 watch deployment 之前先判断 deployment 是否就绪, 如果 deployment 已经
//    // 就绪了,就没必要 watch 了.
//    if h.IsReady(name) {
//        return nil
//    }
//    // 检查 deploymen 是否存在,如果存在则不用 watch
//    if _, err = h.Get(name); err != nil {
//        return err
//    }

//    for {
//        // 开始监听 deployment 事件,
//        // 1. 如果监听到了 modified 事件, 就检查下 deployment 的状态.
//        //    如果 conditions.Type == appsv1.DeploymentAvailable 并且 conditions.Status == corev1.ConditionTrue
//        //    说明 deployment 已经准备好了.
//        // 2. 如果监听到 watch.Deleted 事件, 说明 deployment 已经删除了, 不需要再监听了
//        // 3. 注意这个 watcher 的创建必须在第一层循环之内,如果 watcher 在循环外面
//        //    当 apiserver timeout, 这个 watcher 会被 close, 即使再循环也是在使用
//        //    一个被 close 掉的 watcher, 就无法再监听 deployment.
//        listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
//        listOptions.TimeoutSeconds = &timeout
//        watcher, err = h.clientset.AppsV1().Deployments(h.namespace).Watch(h.ctx, listOptions)
//        if err != nil {
//            return err
//        }
//        // 连接 kubernetes 是有通过 http/https 方式连接的, 有一个 keepalived 的时间
//        // 时间一到, 就会断开 kubernetes 的连接, 此时  watch.ResultChan 通道就会关闭.
//        // 所以说, 这个方法 WaitReady 等待 deployment 处于就绪状态的最长等待时间就是
//        // 连接 kubernetes 的 keepalive 时间. 好像是10分钟
//        for event := range watcher.ResultChan() {
//            switch event.Type {
//            case watch.Modified:
//                if h.IsReady(name) {
//                    watcher.Stop()
//                    return nil
//                }
//            case watch.Deleted: // deployment 已经删除了, 退出监听
//                watcher.Stop()
//                return fmt.Errorf("%s deleted", name)
//            case watch.Bookmark:
//                log.Debug("watch deployment: bookmark.")
//            case watch.Error:
//                log.Debug("watch deployment: error")
//            }
//        }
//        // If event channel is closed, it means the server has closed the connection.
//        log.Debug("watch deployment: reconnect to kubernetes.")
//        watcher.Stop()
//    }
//}

// GetRS get all replicaset created by the deployment.
func (h *Handler) GetRS(object interface{}) ([]*appsv1.ReplicaSet, error) {
	switch val := object.(type) {
	// if object type is string, the object is regarded as deployment name,
	// and check whether deployment exists.
	case string:
		deploy, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getRS(deploy)
	case *appsv1.Deployment:
		return h.getRS(val)
	case appsv1.Deployment:
		return h.getRS(&val)
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getRS(deploy *appsv1.Deployment) ([]*appsv1.ReplicaSet, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	rsList, err := h.clientset.AppsV1().ReplicaSets(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}

	var rl []*appsv1.ReplicaSet
	for i := range rsList.Items {
		for _, or := range rsList.Items[i].OwnerReferences {
			if or.Name == deploy.Name {
				rl = append(rl, &rsList.Items[i])
			}
		}
	}
	return rl, nil
}

// GetPods get all pods created by the deployment.
func (h *Handler) GetPods(object interface{}) ([]*corev1.Pod, error) {
	// GetPods does not need to check deployment is exists.
	// GetRS will check it.
	rsList, err := h.GetRS(object)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	podList, err := h.clientset.CoreV1().Pods(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}

	var pl []*corev1.Pod
	for i := range rsList {
		for j := range podList.Items {
			for _, or := range podList.Items[j].OwnerReferences {
				if or.Name == rsList[i].Name {
					pl = append(pl, &podList.Items[j])
				}
			}
		}
	}
	return pl, nil
}

// GetPVC get all persistentvolumeclaims mounted by the deployment.
func (h *Handler) GetPVC(object interface{}) ([]string, error) {
	switch val := object.(type) {
	// if object type is string, this object is regarded as deployment name,
	// and check whether deployment exists.
	case string:
		deploy, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getPVC(deploy), nil
	case *appsv1.Deployment:
		return h.getPVC(val), nil
	case appsv1.Deployment:
		return h.getPVC(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getPVC(deploy *appsv1.Deployment) []string {
	var pl []string
	for _, volume := range deploy.Spec.Template.Spec.Volumes {
		// 有些 volume.PersistentVolumeClaim 是不存在的, 其值默认是 nil 如果不加以判断就直接获取
		// volume.PersistentVolumeClaim.ClaimName, 就操作了非法地址, 程序会直接 panic.
		if volume.PersistentVolumeClaim != nil {
			pl = append(pl, volume.PersistentVolumeClaim.ClaimName)
		}
	}
	return pl
}

// GetPV get all persistentvolumes mounted by the deployment.
func (h *Handler) GetPV(object interface{}) ([]string, error) {
	// GetPV does not need to check whether deployment is exists.
	// GetPVC will do it.
	pvcList, err := h.GetPVC(object)
	if err != nil {
		return nil, err
	}

	// pvc.spec.volumeName 的值就是 pvc 中绑定的 pv 的名字
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

// GetAge returns deployment age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		deploy, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		return time.Now().Sub(deploy.CreationTimestamp.Time), nil
	case *appsv1.Deployment:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case appsv1.Deployment:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}

// GetContainers get all container of this deployment.
func (h *Handler) GetContainers(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		sts, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getContainers(sts), nil
	case *appsv1.Deployment:
		return h.getContainers(val), nil
	case appsv1.Deployment:
		return h.getContainers(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getContainers(sts *appsv1.Deployment) []string {
	var cl []string
	for _, container := range sts.Spec.Template.Spec.Containers {
		cl = append(cl, container.Name)
	}
	return cl
}

// GetImages get all container images of this deployment.
func (h *Handler) GetImages(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		sts, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getImages(sts), nil
	case *appsv1.Deployment:
		return h.getImages(val), nil
	case appsv1.Deployment:
		return h.getImages(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getImages(sts *appsv1.Deployment) []string {
	var il []string
	for _, container := range sts.Spec.Template.Spec.Containers {
		il = append(il, container.Image)
	}
	return il
}
