package daemonset

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

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

// WaitReady wait the daemonset to be th ready status.
func (h *Handler) WaitReady(name string) (err error) {
	var (
		watcher watch.Interface
		timeout = int64(0)
	)
	// 如果 daemonset 已经就绪, 就没必要继续 watch 了
	if h.IsReady(name) {
		return nil
	}
	// 判断 daemonset 是否存在
	if _, err = h.Get(name); err != nil {
		return err
	}
	for {
		listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
		listOptions.TimeoutSeconds = &timeout
		watcher, err = h.clientset.AppsV1().DaemonSets(h.namespace).Watch(h.ctx, listOptions)
		if err != nil {
			return
		}
		for event := range watcher.ResultChan() {
			switch event.Type {
			case watch.Modified:
				if h.IsReady(name) {
					watcher.Stop()
					return nil
				}
			case watch.Deleted:
				watcher.Stop()
				return fmt.Errorf("%s deleted", name)
			case watch.Bookmark:
				log.Debug("watch daemonset: bookmark.")
			case watch.Error:
				log.Debug("watch daemonset: error")

			}
		}
		log.Debug("watch daemonset: reconnect to kubernetes")
		watcher.Stop()
	}
}

// GetPods all pods created by the daemonset.
func (h *Handler) GetPods(name string) ([]corev1.Pod, error) {
	if _, err := h.Get(name); err != nil {
		return nil, err
	}

	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	// get all pods in the namespace that the daemonset is running.
	podList, err := h.clientset.CoreV1().Pods(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}

	var pl []corev1.Pod
	for _, p := range podList.Items {
		for _, or := range p.OwnerReferences {
			if or.Name == name {
				pl = append(pl, p)
			}
		}
	}
	return pl, nil
}

// GetPVC get daemonset all pvc by name.
func (h *Handler) GetPVC(name string) ([]string, error) {
	daemonset, err := h.Get(name)
	if err != nil {
		return nil, err
	}

	var pl []string
	for _, volume := range daemonset.Spec.Template.Spec.Volumes {
		// 有些 volume.PersistentVolumeClaim 是不存在的, 其值默认是 nil 如果不加以判断就直接获取
		// volume.PersistentVolumeClaim.ClaimName, 就操作了非法地址, 程序会直接 panic.
		if volume.PersistentVolumeClaim != nil {
			pl = append(pl, volume.PersistentVolumeClaim.ClaimName)
		}
	}
	return pl, nil
}

// GetPV get daemonset all pv by name.
func (h *Handler) GetPV(name string) ([]string, error) {
	if _, err := h.Get(name); err != nil {
		return nil, err
	}

	pvcList, err := h.GetPVC(name)
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

// GetAge get daemonset age.
func (h *Handler) GetAge(name string) (time.Duration, error) {
	ds, err := h.Get(name)
	if err != nil {
		return time.Duration(int64(0)), nil
	}

	ctime := ds.CreationTimestamp.Time
	return time.Now().Sub(ctime), nil
}
