package statefulset

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

// IsReady check if the statefulset is ready.
// ref: https://github.com/kubernetes/kubernetes/blob/a1128e380c2cf1c2d7443694673d9f1dd63eb518/staging/src/k8s.io/kubectl/pkg/polymorphichelpers/rollout_status.go#L120
func (h *Handler) IsReady(name string) bool {
	checkGeneration := func(sts *appsv1.StatefulSet) bool {
		if sts.Generation != sts.Status.ObservedGeneration {
			return false
		}
		return true
	}
	checkReplicas := func(sts *appsv1.StatefulSet) bool {
		if sts.Spec.Replicas == nil {
			return false
		}
		if *sts.Spec.Replicas != sts.Status.ReadyReplicas {
			return false
		}
		if *sts.Spec.Replicas != sts.Status.AvailableReplicas {
			return false
		}
		if *sts.Spec.Replicas != sts.Status.CurrentReplicas {
			return false
		}
		if *sts.Spec.Replicas != sts.Status.UpdatedReplicas {
			return false
		}
		if *sts.Spec.Replicas != sts.Status.Replicas {
			return false
		}
		return true
	}
	checkRevision := func(sts *appsv1.StatefulSet) bool {
		if sts.Status.UpdateRevision != sts.Status.CurrentRevision {
			return false
		}
		return true
	}

	sts, err := h.Get(name)
	if err != nil {
		return false
	}
	if checkGeneration(sts) && checkReplicas(sts) && checkRevision(sts) {
		return true
	}
	return false
}

// WaitReady waiting for the statefulset to be in the ready status.
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

	// this goroutine used to check whether statefulset is ready and exists.
	// if statefulset already ready, return nil.
	// if statefulset does not exist, return error.
	// if statefulset exists but not ready, return nothing.
	go func(context.Context) {
		<-chkCh
		for i := 0; i < 6; i++ {
			// statefulset is already ready, return nil
			if h.IsReady(name) {
				errCh <- nil
				return
			}
			// statefulset no longer exists, return err
			_, err := h.Get(name)
			if k8serrors.IsNotFound(err) {
				errCh <- err
				return
			}
			// statefulset exists but not ready, return nothing.
			if err == nil {
				return
			}
			log.Error(err)
			time.Sleep(time.Second * 10)
		}
	}(ctxCheck)

	// this goroutine used to watch statefulset.
	go func(ctx context.Context) {
		for {
			timeout := int64(0)
			listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
			listOptions.TimeoutSeconds = &timeout
			watcher, err := h.clientset.AppsV1().StatefulSets(h.namespace).Watch(h.ctx, listOptions)
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
					errCh <- fmt.Errorf("statefulset/%s was deleted", name)
					return
				case watch.Bookmark:
					log.Debug("watch statefulset: bookmark")
				case watch.Error:
					log.Debug("watch statefulset: error")
				}
			}
			// If event channel is closed, it means the kube-apiserver has closed the connection.
			log.Debug("watch statefulset: reconnect to kubernetes")
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

//// WaitReady wait the statefulset to be in the ready status.
//func (h *Handler) WaitReady(name string) (err error) {
//    var (
//        watcher watch.Interface
//        timeout = int64(0)
//    )
//    // watch 之前先判断 statefulset 是否就绪, 如果已经继续就没必要继续 watch 了
//    if h.IsReady(name) {
//        return nil
//    }
//    // 是否判断 statefulset 是否存在
//    if _, err = h.Get(name); err != nil {
//        return err
//    }

//    // 1.由于 watcher 会因为 keepalive 超时被 kube-apiserver 中断, 所以需要循环创建 watcher
//    // 2.这个 watcher 要放在第一层 for 循环里面
//    for {
//        listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
//        listOptions.TimeoutSeconds = &timeout
//        watcher, err = h.clientset.AppsV1().StatefulSets(h.namespace).Watch(h.ctx, listOptions)
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
//                log.Debug("watch statefulset: bookmark")
//            case watch.Error:
//                log.Debug("watch statefulset: error")
//            }
//        }
//        log.Debug("watch statefulset: reconnect to kubernetes")
//        watcher.Stop()
//    }
//}

// GetPods get all pods created by the statefulset.
func (h *Handler) GetPods(object interface{}) ([]*corev1.Pod, error) {
	// if statefulset not exist, return err.
	switch val := object.(type) {
	case string:
		sts, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getPods(sts)
	case *appsv1.StatefulSet:
		return h.getPods(val)
	case appsv1.StatefulSet:
		return h.getPods(&val)
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getPods(sts *appsv1.StatefulSet) ([]*corev1.Pod, error) {
	// get all pods in the namespace that statefulset is running.
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	podList, err := h.clientset.CoreV1().Pods(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}

	var pl []*corev1.Pod
	for i := range podList.Items {
		for _, or := range podList.Items[i].OwnerReferences {
			if or.Name == sts.Name {
				pl = append(pl, &podList.Items[i])
			}
		}
	}
	return pl, nil
}

// GetPVC get all persistentvolumeclaims mounted by the statefulset.
func (h *Handler) GetPVC(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		sts, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getPVC(sts), nil
	case *appsv1.StatefulSet:
		return h.getPVC(val), nil
	case appsv1.StatefulSet:
		return h.getPVC(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getPVC(sts *appsv1.StatefulSet) []string {
	// pvc 的格式为 sts name + VolumeClaimTemplates name + replicas 编号
	var pl []string
	for _, p := range sts.Spec.VolumeClaimTemplates {
		for i := int32(0); i < *sts.Spec.Replicas; i++ {
			pl = append(pl, fmt.Sprintf("%s-%s-%d", p.ObjectMeta.Name, sts.Name, i))
		}
	}
	return pl
}

// GetPV get all persistentvolumes mounted by the statefulset.
func (h *Handler) GetPV(object interface{}) ([]string, error) {
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

// GetAge returns the statefulset age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		sts, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		ctime := sts.CreationTimestamp.Time
		return time.Now().Sub(ctime), nil
	case *appsv1.StatefulSet:
		ctime := val.CreationTimestamp.Time
		return time.Now().Sub(ctime), nil
	case appsv1.StatefulSet:
		ctime := val.CreationTimestamp.Time
		return time.Now().Sub(ctime), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}

// GetContainers get all container of this statefulset.
func (h *Handler) GetContainers(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		sts, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getContainers(sts), nil
	case *appsv1.StatefulSet:
		return h.getContainers(val), nil
	case appsv1.StatefulSet:
		return h.getContainers(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getContainers(sts *appsv1.StatefulSet) []string {
	var cl []string
	for _, container := range sts.Spec.Template.Spec.Containers {
		cl = append(cl, container.Name)
	}
	return cl
}

// GetImages get all container images of this statefulset.
func (h *Handler) GetImages(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		sts, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getImages(sts), nil
	case *appsv1.StatefulSet:
		return h.getImages(val), nil
	case appsv1.StatefulSet:
		return h.getImages(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getImages(sts *appsv1.StatefulSet) []string {
	var il []string
	for _, container := range sts.Spec.Template.Spec.Containers {
		il = append(il, container.Image)
	}
	return il
}
