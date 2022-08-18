package daemonset

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

var ERR_TYPE = fmt.Errorf("type must be *appsv1.DaemonSet, appsv1.DaemonSet or string")

// IsReady check if the daemonset is ready.
func (h *Handler) IsReady(name string) bool {
	daemonset, err := h.Get(name)
	if err != nil {
		return false
	}
	//log.SetLevel(log.TraceLevel)
	//log.Debug(daemonset.Status)
	//log.Debug(daemonset.Status.DesiredNumberScheduled)
	//log.Debug(daemonset.Status.CurrentNumberScheduled)
	//log.Debug(daemonset.Status.NumberAvailable)
	//log.Debug(daemonset.Status.NumberReady)
	desiredNumberScheduled := daemonset.Status.DesiredNumberScheduled
	// if desiredNumberScheduled = 0, that means daemonset not ready.
	if desiredNumberScheduled == 0 {
		return false
	}
	if daemonset.Status.CurrentNumberScheduled == desiredNumberScheduled &&
		daemonset.Status.NumberAvailable == desiredNumberScheduled &&
		daemonset.Status.NumberReady == desiredNumberScheduled {
		return true
	}
	return false
}

// WaitReady waiting for the daemonset to be in the ready status.
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

	// the goroutine used to check whether daemonset is ready and exists.
	// if daemonset already ready, return nil.
	// if daemonset does not exist, return error.
	// if daemonset exists but not ready, return nothing.
	go func(context.Context) {
		<-chkCh
		for i := 0; i < 6; i++ {
			// daemonset is already ready, return nil
			if h.IsReady(name) {
				errCh <- nil
				return
			}
			// daemonset no longer exists, return err
			_, err := h.Get(name)
			if k8serrors.IsNotFound(err) {
				errCh <- err
				return
			}
			// daemonset exists but not ready, return nothing.
			if err == nil {
				return
			}
			log.Error(err)
			time.Sleep(time.Second * 10)
		}
	}(ctxCheck)

	// this goroutine used to watch daemonset.
	go func(ctx context.Context) {
		for {
			timeout := int64(0)
			listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
			listOptions.TimeoutSeconds = &timeout
			watcher, err := h.clientset.AppsV1().DaemonSets(h.namespace).Watch(h.ctx, listOptions)
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
					errCh <- fmt.Errorf("daemonset/%s was deleted", name)
					return
				case watch.Bookmark:
					log.Debug("watch daemonset: bookmark")
				case watch.Error:
					log.Debug("watch daemonset: error")
				}
			}
			// If event channel is closed, it means the kube-apiserver has closed the connection.
			log.Debug("watch daemonset: reconnect to kubernetes")
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

//// WaitReady wait the daemonset to be th ready status.
//func (h *Handler) WaitReady2(name string) (err error) {
//    var (
//        watcher watch.Interface
//        timeout = int64(0)
//    )
//    // 如果 daemonset 已经就绪, 就没必要继续 watch 了
//    if h.IsReady(name) {
//        return nil
//    }
//    // 判断 daemonset 是否存在
//    if _, err = h.Get(name); err != nil {
//        return err
//    }
//    for {
//        listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
//        listOptions.TimeoutSeconds = &timeout
//        watcher, err = h.clientset.AppsV1().DaemonSets(h.namespace).Watch(h.ctx, listOptions)
//        if err != nil {
//            return
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
//                log.Debug("watch daemonset: bookmark.")
//            case watch.Error:
//                log.Debug("watch daemonset: error")

//            }
//        }
//        log.Debug("watch daemonset: reconnect to kubernetes")
//        watcher.Stop()
//    }
//}

// GetPods get all pods created by the daemonset.
func (h *Handler) GetPods(object interface{}) ([]*corev1.Pod, error) {
	switch val := object.(type) {
	case string:
		ds, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getPods(ds)
	case *appsv1.DaemonSet:
		return h.getPods(val)
	case appsv1.DaemonSet:
		return h.getPods(&val)
	default:
		return nil, ERR_TYPE
	}
}
func (h *Handler) getPods(ds *appsv1.DaemonSet) ([]*corev1.Pod, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	// get all pods in the namespace that the daemonset is running.
	podList, err := h.clientset.CoreV1().Pods(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}

	var pl []*corev1.Pod
	for i := range podList.Items {
		for _, or := range podList.Items[i].OwnerReferences {
			if or.Name == ds.Name {
				pl = append(pl, &podList.Items[i])
			}
		}
	}
	return pl, nil
}

// GetPVC get all persistentvolumeclaims mounted by the daemonset.
func (h *Handler) GetPVC(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		ds, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getPVC(ds), nil
	case *appsv1.DaemonSet:
		return h.getPVC(val), nil
	case appsv1.DaemonSet:
		return h.getPVC(&val), nil
	default:
		return nil, ERR_TYPE
	}

}
func (h *Handler) getPVC(ds *appsv1.DaemonSet) []string {
	var pl []string
	for _, volume := range ds.Spec.Template.Spec.Volumes {
		// 有些 volume.PersistentVolumeClaim 是不存在的, 其值默认是 nil 如果不加以判断就直接获取
		// volume.PersistentVolumeClaim.ClaimName, 就操作了非法地址, 程序会直接 panic.
		if volume.PersistentVolumeClaim != nil {
			pl = append(pl, volume.PersistentVolumeClaim.ClaimName)
		}
	}
	return pl
}

// GetPV get all persistentvolumes mounted by the daemonset.
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

// NumDesired returns the total number of pods that should be running.
func (h *Handler) NumDesired(object interface{}) (int32, error) {
	switch val := object.(type) {
	case string:
		ds, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return ds.Status.DesiredNumberScheduled, nil
	case *appsv1.DaemonSet:
		return val.Status.DesiredNumberScheduled, nil
	case appsv1.DaemonSet:
		return val.Status.DesiredNumberScheduled, nil
	default:
		return 0, ERR_TYPE
	}
}

// NumCurrent returns the total number of pods that currently running.
func (h *Handler) NumCurrent(object interface{}) (int32, error) {
	switch val := object.(type) {
	case string:
		ds, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return ds.Status.CurrentNumberScheduled, nil
	case *appsv1.DaemonSet:
		return val.Status.CurrentNumberScheduled, nil
	case appsv1.DaemonSet:
		return val.Status.CurrentNumberScheduled, nil
	default:
		return 0, ERR_TYPE
	}
}

// NumReady
func (h *Handler) NumReady(object interface{}) (int32, error) {
	switch val := object.(type) {
	case string:
		ds, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return ds.Status.NumberReady, nil
	case *appsv1.DaemonSet:
		return val.Status.NumberReady, nil
	case appsv1.DaemonSet:
		return val.Status.NumberReady, nil
	default:
		return 0, ERR_TYPE
	}
}

// NumAvailable
func (h *Handler) NumAvailable(object interface{}) (int32, error) {
	switch val := object.(type) {
	case string:
		ds, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return ds.Status.NumberAvailable, nil
	case *appsv1.DaemonSet:
		return val.Status.NumberAvailable, nil
	case appsv1.DaemonSet:
		return val.Status.NumberAvailable, nil
	default:
		return 0, ERR_TYPE
	}
}

// GetAge returns daemonset age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		ds, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), nil
		}
		return time.Now().Sub(ds.CreationTimestamp.Time), nil
	case *appsv1.DaemonSet:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case appsv1.DaemonSet:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(int64(0)), ERR_TYPE
	}
}

// GetContainers get all container of this daemonset.
func (h *Handler) GetContainers(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		sts, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getContainers(sts), nil
	case *appsv1.DaemonSet:
		return h.getContainers(val), nil
	case appsv1.DaemonSet:
		return h.getContainers(&val), nil
	default:
		return nil, ERR_TYPE
	}
}
func (h *Handler) getContainers(sts *appsv1.DaemonSet) []string {
	var cl []string
	for _, container := range sts.Spec.Template.Spec.Containers {
		cl = append(cl, container.Name)
	}
	return cl
}

// GetImages get all container images of this daemonset.
func (h *Handler) GetImages(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		sts, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getImages(sts), nil
	case *appsv1.DaemonSet:
		return h.getImages(val), nil
	case appsv1.DaemonSet:
		return h.getImages(&val), nil
	default:
		return nil, ERR_TYPE
	}
}
func (h *Handler) getImages(sts *appsv1.DaemonSet) []string {
	var il []string
	for _, container := range sts.Spec.Template.Spec.Containers {
		il = append(il, container.Image)
	}
	return il
}
