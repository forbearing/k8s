package replicaset

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

// IsReady check if the replicaset is ready
func (h *Handler) IsReady(name string) bool {
	checkGeneration := func(rs *appsv1.ReplicaSet) bool {
		if rs.Generation != rs.Status.ObservedGeneration {
			return false
		}
		return true
	}
	checkReplicas := func(rs *appsv1.ReplicaSet) bool {
		if rs.Spec.Replicas == nil {
			return false
		}
		if *rs.Spec.Replicas != rs.Status.ReadyReplicas {
			return false
		}
		if *rs.Spec.Replicas != rs.Status.AvailableReplicas {
			return false
		}
		if *rs.Spec.Replicas != rs.Status.Replicas {
			return false
		}
		if *rs.Spec.Replicas != rs.Status.FullyLabeledReplicas {

		}
		return true
	}

	rs, err := h.Get(name)
	if err != nil {
		return false
	}
	if checkGeneration(rs) && checkReplicas(rs) {
		return true
	}
	return false
}

// WaitReady waiting for the replicaset to be in the ready status.
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

	// this goroutine used to check whether replicaset is ready and exists.
	// if replicaset already ready, return nil.
	// if replicaset does not exist, return error.
	// if replicaset exists but not ready, return nothing.
	go func(context.Context) {
		<-chkCh
		for i := 0; i < 6; i++ {
			// replicaset is already ready, return nil
			if h.IsReady(name) {
				errCh <- nil
				return
			}
			// replicaset no longer exists, return err
			_, err := h.Get(name)
			if k8serrors.IsNotFound(err) {
				errCh <- err
				return
			}
			// replicaset exists but not ready, return nothing.
			if err == nil {
				return
			}
			log.Error(err)
			time.Sleep(time.Second * 10)
		}
	}(ctxCheck)

	// this goroutine used to watch replicaset.
	go func(ctx context.Context) {
		for {
			timeout := int64(0)
			listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
			listOptions.TimeoutSeconds = &timeout
			watcher, err := h.clientset.AppsV1().ReplicaSets(h.namespace).Watch(h.ctx, listOptions)
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
					errCh <- fmt.Errorf("replicaset/%s was deleted", name)
					return
				case watch.Bookmark:
					log.Debug("watch replicaset: bookmark")
				case watch.Error:
					log.Debug("watch replicaset: error")
				}
			}
			// If event channel is closed, it means the kube-apiserver has closed the connection.
			log.Debug("watch replicaset: reconnect to kubernetes")
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

//// WaitReady wait the replicaset to be th ready status
//func (h *Handler) WaitReady(name string, check bool) (err error) {
//    var (
//        watcher watch.Interface
//        timeout = int64(0)
//    )
//    // 如果 replicaset 已经就绪, 就没必要继续 watch 了
//    if h.IsReady(name) {
//        return
//    }
//    // 是否判断 replicaset 是否存在
//    if check {
//        if _, err = h.Get(name); err != nil {
//            return
//        }
//    }
//    for {
//        listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
//        listOptions.TimeoutSeconds = &timeout
//        watcher, err = h.clientset.AppsV1().ReplicaSets(h.namespace).Watch(h.ctx, listOptions)
//        if err != nil {
//            return
//        }
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
//                log.Debug("watch replicaset: bookmark.")
//            case watch.Error:
//                log.Debug("watch replicaset: error")

//            }
//        }
//        log.Debug("watch replicaset: reconnect to kubernetes")
//        watcher.Stop()
//    }
//}

// GetPods get replicaset all pods
func (h *Handler) GetPods(object interface{}) ([]*corev1.Pod, error) {
	switch val := object.(type) {
	case string:
		rs, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getPods(rs)
	case *appsv1.ReplicaSet:
		return h.getPods(val)
	case appsv1.ReplicaSet:
		return h.getPods(&val)
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getPods(rs *appsv1.ReplicaSet) ([]*corev1.Pod, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	podList, err := h.clientset.CoreV1().Pods(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}

	var pl []*corev1.Pod
	for i := range podList.Items {
		for _, or := range podList.Items[i].OwnerReferences {
			if or.Name == rs.Name {
				pl = append(pl, &podList.Items[i])
			}
		}
	}
	return pl, nil
}

// GetPVC get all persistentvolumeclaims mounted by the replicaset.
func (h *Handler) GetPVC(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		rs, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getPVC(rs), nil
	case *appsv1.ReplicaSet:
		return h.getPVC(val), nil
	case appsv1.ReplicaSet:
		return h.getPVC(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getPVC(rs *appsv1.ReplicaSet) []string {
	var pl []string
	for _, volume := range rs.Spec.Template.Spec.Volumes {
		if volume.PersistentVolumeClaim != nil {
			pl = append(pl, volume.PersistentVolumeClaim.ClaimName)
		}
	}
	return pl
}

// GetPV get all persistentvolumes mounted by the replicaset.
func (h *Handler) GetPV(object interface{}) ([]string, error) {
	// GetPV does not need to check whether replicaset is exists.
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

// GetAge returns replicaset age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		rs, err := h.Get(val)
		if err != nil {
			return time.Duration(0), err
		}
		return time.Now().Sub(rs.CreationTimestamp.Time), nil
	case *appsv1.ReplicaSet:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case appsv1.ReplicaSet:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(0), ErrInvalidToolsType
	}
}
