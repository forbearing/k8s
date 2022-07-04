package replicaset

import (
	"fmt"
	"strings"

	"github.com/forbearing/k8s/persistentvolumeclaim"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// GetPods get replicaset all pods
func (h *Handler) GetPods(name string) (podList []string, err error) {
	// 检查 replicaset 是否就绪
	err = h.WaitReady(name, true)
	if err != nil {
		return
	}
	if !h.IsReady(name) {
		err = fmt.Errorf("replicaset %s not ready", name)
		return
	}

	// 创建一个 appsv1.ReplicaSet 对象
	replicaset, err := h.Get(name)
	if err != nil {
		return
	}
	// 通过 spec.selector.matchLabels 找到 replicaset 创建的 pod
	matchLabels := replicaset.Spec.Selector.MatchLabels
	labelSelector := ""
	for key, value := range matchLabels {
		labelSelector = labelSelector + fmt.Sprintf("%s=%s,", key, value)
	}
	labelSelector = strings.TrimRight(labelSelector, ",")
	podObjList, err := h.clientset.CoreV1().Pods(h.namespace).List(h.ctx,
		metav1.ListOptions{LabelSelector: labelSelector})
	for _, pod := range podObjList.Items {
		podList = append(podList, pod.Name)
	}
	return
}

// GetPV get replicaset pv by name
func (h *Handler) GetPV(name string) (pvList []string, err error) {
	var (
		pvcHandler *persistentvolumeclaim.Handler
		pvcObj     *corev1.PersistentVolumeClaim
		pvcList    []string
	)
	err = h.WaitReady(name, true)
	if err != nil {
		return
	}
	if !h.IsReady(name) {
		err = fmt.Errorf("replicaset %s not ready", name)
		return
	}

	pvcHandler, err = persistentvolumeclaim.New(h.ctx, h.namespace, h.kubeconfig)
	if err != nil {
		return
	}
	pvcList, err = h.GetPVC(name)
	if err != nil {
		return
	}

	for _, pvcName := range pvcList {
		pvcObj, err = pvcHandler.Get(pvcName)
		if err != nil {
			return
		}
		pvList = append(pvList, pvcObj.Spec.VolumeName)
	}

	return
}

// GetPVC get replicaset pvc by name
func (h *Handler) GetPVC(name string) (pvcList []string, err error) {
	err = h.WaitReady(name, true)
	if err != nil {
		return
	}
	if !h.IsReady(name) {
		err = fmt.Errorf("replicaset %s not ready", name)
		return
	}
	replicaset, err := h.Get(name)
	if err != nil {
		return
	}
	// 如果 volume.PersistentVolumeClaim 为 nil, 就不能再操作 volume.PersistentVolumeClaim.ClaimName
	for _, volume := range replicaset.Spec.Template.Spec.Volumes {
		if volume.PersistentVolumeClaim != nil {
			pvcList = append(pvcList, volume.PersistentVolumeClaim.ClaimName)
		}
	}
	return
}

// IsReady check if the replicaset is ready
func (h *Handler) IsReady(name string) bool {
	replicaset, err := h.Get(name)
	if err != nil {
		return false
	}
	replicas := replicaset.Status.Replicas
	if replicaset.Status.AvailableReplicas == replicas &&
		replicaset.Status.FullyLabeledReplicas == replicas &&
		replicaset.Status.ReadyReplicas == replicas {
		return true
	}
	return false
}

// WaitReady wait the replicaset to be th ready status
func (h *Handler) WaitReady(name string, check bool) (err error) {
	var (
		watcher watch.Interface
		timeout = int64(0)
	)
	// 如果 replicaset 已经就绪, 就没必要继续 watch 了
	if h.IsReady(name) {
		return
	}
	// 是否判断 replicaset 是否存在
	if check {
		if _, err = h.Get(name); err != nil {
			return
		}
	}
	for {
		listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
		listOptions.TimeoutSeconds = &timeout
		watcher, err = h.clientset.AppsV1().ReplicaSets(h.namespace).Watch(h.ctx, listOptions)
		if err != nil {
			return
		}
		for event := range watcher.ResultChan() {
			switch event.Type {
			case watch.Modified:
				if h.IsReady(name) {
					watcher.Stop()
					return
				}
			case watch.Deleted:
				watcher.Stop()
				return fmt.Errorf("%s deleted", name)
			case watch.Bookmark:
				log.Debug("watch replicaset: bookmark.")
			case watch.Error:
				log.Debug("watch replicaset: error")

			}
		}
		log.Debug("watch replicaset: reconnect to kubernetes")
		watcher.Stop()
	}
}
