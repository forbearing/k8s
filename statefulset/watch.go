package statefulset

import (
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName watch statefulset by name.
func (h *Handler) WatchByName(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher watch.Interface
		timeout = int64(0)
		isExist bool
	)
	for {
		listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
		listOptions.TimeoutSeconds = &timeout
		if watcher, err = h.clientset.AppsV1().StatefulSets(h.namespace).Watch(h.ctx, listOptions); err != nil {
			return
		}
		if _, err = h.Get(name); err != nil {
			isExist = false // statefulset not exist
		} else {
			isExist = true // statefulset exist
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
				log.Debug("watch statefulset: bookmark")
			case watch.Error:
				log.Debug("watch statefulset: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch statefulset: reconnect to kubernetes")
		watcher.Stop()
	}
}

// WatchByLabel watch statefulset by label.
func (h *Handler) WatchByLabel(labelSelector string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher         watch.Interface
		statefulsetList *appsv1.StatefulSetList
		timeout         = int64(0)
		isExist         bool
	)
	for {
		if watcher, err = h.clientset.AppsV1().StatefulSets(h.namespace).Watch(h.ctx,
			metav1.ListOptions{LabelSelector: labelSelector, TimeoutSeconds: &timeout}); err != nil {
			return
		}
		if statefulsetList, err = h.List(labelSelector); err != nil {
			return
		}
		if len(statefulsetList.Items) == 0 {
			isExist = false // statefulset not exist
		} else {
			isExist = true // statefulset exist
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
				log.Debug("watch statefulset: bookmark")
			case watch.Error:
				log.Debug("watch statefulset: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch statefulset: reconnect to kubernetes")
		watcher.Stop()
	}
}

// Watch watch statefulset by name, alias to "WatchByName".
func (h *Handler) Watch(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	return h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
}
