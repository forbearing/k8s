package replicationcontroller

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// IsReady check if the replicationcontroller is ready
func (h *Handler) IsReady(name string) bool {
	// 获取 *corev1.ReplicationController 对象
	rc, err := h.Get(name)
	if err != nil {
		return false
	}
	replicas := rc.Status.Replicas
	if rc.Status.AvailableReplicas == replicas &&
		rc.Status.FullyLabeledReplicas == replicas &&
		rc.Status.ReadyReplicas == replicas {
		return true
	}
	return false
}

// WaitReady wait for the replicationcontroller to be in the ready status
func (h *Handler) WaitReady(name string, check bool) (err error) {
	var (
		watcher watch.Interface
		timeout = int64(0)
	)
	// 在 watch 之前就先判断 replicationcontroller 是否就绪, 如果就绪了就没必要 watch 了
	if h.IsReady(name) {
		return
	}
	// 是否判断 replicationcontroller 是否存在
	if check {
		if _, err = h.Get(name); err != nil {
			return
		}
	}
	for {
		// replicationcontroller 没有就绪, 那么就开始监听 replicationcontroller 的事件
		listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
		listOptions.TimeoutSeconds = &timeout
		watcher, err = h.clientset.CoreV1().ReplicationControllers(h.namespace).Watch(h.ctx, listOptions)
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
				log.Debug("watch replicationcontroller: bookmark")
			case watch.Error:
				log.Debug("watch replicationcontroller: error")
			}
		}
		// watcher 因为 keepalive 超时断开了连接, 关闭了 channel
		log.Debug("watch replicationcontroller: reconnect to kubernetes")
		watcher.Stop()
	}
}
