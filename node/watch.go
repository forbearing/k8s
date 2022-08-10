package node

import (
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName watch node by name.
func (h *Handler) WatchByName(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher watch.Interface
		timeout = int64(0)
		isExist bool
	)
	for {
		listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name})
		listOptions.TimeoutSeconds = &timeout
		if watcher, err = h.clientset.CoreV1().Nodes().Watch(h.ctx, listOptions); err != nil {
			return
		}
		if _, err = h.Get(name); err != nil {
			isExist = false // node not exist
		} else {
			isExist = true // node exist
		}
		for event := range watcher.ResultChan() {
			switch event.Type {
			case watch.Added:
				if !isExist {
					addFunc(x)
				}
				isExist = true
			case watch.Modified:
				modifyFunc(x)
				isExist = true
			case watch.Deleted:
				deleteFunc(x)
				isExist = false
			case watch.Bookmark:
				log.Debug("watch node: bookmark")
			case watch.Error:
				log.Debug("watch node: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch node: reconnect to kubernetes")
		watcher.Stop()
	}
}

// WatchByLabel watch node by labels.
func (h *Handler) WatchByLabel(labels string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher  watch.Interface
		nodeList []*corev1.Node
		timeout  = int64(0)
		isExist  bool
	)
	for {
		if watcher, err = h.clientset.CoreV1().Nodes().Watch(h.ctx,
			metav1.ListOptions{LabelSelector: labels, TimeoutSeconds: &timeout}); err != nil {
			return
		}
		if nodeList, err = h.ListByLabel(labels); err != nil {
			return
		}
		if len(nodeList) == 0 {
			isExist = false // node not exist
		} else {
			isExist = true // node exist
		}
		for event := range watcher.ResultChan() {
			switch event.Type {
			case watch.Added:
				if !isExist {
					addFunc(x)
				}
				isExist = true
			case watch.Modified:
				modifyFunc(x)
				isExist = true
			case watch.Deleted:
				deleteFunc(x)
				isExist = false
			case watch.Bookmark:
				log.Debug("watch node: bookmark")
			case watch.Error:
				log.Debug("watch node: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch node: reconnect to kubernetes")
		watcher.Stop()
	}
}

// WatchByName watch node by name, alias to "WatchByName".
func (h *Handler) Watch(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	return h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
}
