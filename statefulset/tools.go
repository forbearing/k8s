package statefulset

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// IsReady check if the statefulset is ready.
func (h *Handler) IsReady(name string) bool {
	statefulset, err := h.Get(name)
	if err != nil {
		return false
	}

	// 如果 statefulset 的 replicaas 等于 status.AvailableReplicas 的个数
	// 就表明 statefulset 的所有 pod 都就绪了.
	if *statefulset.Spec.Replicas == statefulset.Status.AvailableReplicas {
		return true
	}

	return false
}

// WaitReady wait the statefulset to be in the ready status.
func (h *Handler) WaitReady(name string) (err error) {
	var (
		watcher watch.Interface
		timeout = int64(0)
	)
	// watch 之前先判断 statefulset 是否就绪, 如果已经继续就没必要继续 watch 了
	if h.IsReady(name) {
		return nil
	}
	// 是否判断 statefulset 是否存在
	if _, err = h.Get(name); err != nil {
		return err
	}

	// 1.由于 watcher 会因为 keepalive 超时被 kube-apiserver 中断, 所以需要循环创建 watcher
	// 2.这个 watcher 要放在第一层 for 循环里面
	for {
		listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
		listOptions.TimeoutSeconds = &timeout
		watcher, err = h.clientset.AppsV1().StatefulSets(h.namespace).Watch(h.ctx, listOptions)
		if err != nil {
			return err
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
				log.Debug("watch statefulset: bookmark")
			case watch.Error:
				log.Debug("watch statefulset: error")
			}
		}
		log.Debug("watch statefulset: reconnect to kubernetes")
		watcher.Stop()
	}
}

// GetPods all pods created by the statefulset.
func (h *Handler) GetPods(name string) ([]corev1.Pod, error) {
	// if statefulset not exist, return err.
	if _, err := h.Get(name); err != nil {
		return nil, err
	}

	// get all pods in the namespace that statefulset is running.
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	podList, err := h.clientset.CoreV1().Pods(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}

	var pl []corev1.Pod
	for _, pod := range podList.Items {
		for _, or := range pod.OwnerReferences {
			if or.Name == name {
				pl = append(pl, pod)
			}
		}
	}
	return pl, nil
}

// GetPVC get statefulset all pvc by name.
func (h *Handler) GetPVC(name string) ([]string, error) {
	sts, err := h.Get(name)
	if err != nil {
		return nil, err
	}

	// pvc 的格式为 sts name + VolumeClaimTemplates name + replicas 编号
	var pl []string
	for _, p := range sts.Spec.VolumeClaimTemplates {
		for i := int32(0); i < *sts.Spec.Replicas; i++ {
			pl = append(pl, fmt.Sprintf("%s-%s-%d", p.ObjectMeta.Name, sts.Name, i))
		}
	}
	return pl, nil
}

// GetPV get statefulset all pv by name.
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

// GetAge get statefulset age.
func (h *Handler) GetAge(name string) (time.Duration, error) {
	sts, err := h.Get(name)
	if err != nil {
		return time.Duration(int64(0)), err
	}

	ctime := sts.CreationTimestamp.Time
	return time.Now().Sub(ctime), nil
}
