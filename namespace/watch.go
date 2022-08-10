package namespace

import (
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName watch namespace by name.
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
		if watcher, err = h.clientset.CoreV1().Namespaces().Watch(h.ctx, listOptions); err != nil {
			return
		}
		if _, err = h.Get(name); err != nil {
			isExist = false // namespace not exist bool
		} else {
			isExist = true // namespace exist
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
				log.Debug("watch namespace: bookmark.")
			case watch.Error:
				log.Debug("watch namespace: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch namespace: reconnect to kubernetes")
	}
}

// WatchByLabel watch namespace by labels.
func (h *Handler) WatchByLabel(labels string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher watch.Interface
		nsList  []*corev1.Namespace
		timeout = int64(0)
		isExist bool
	)
	for {
		if watcher, err = h.clientset.CoreV1().Namespaces().Watch(h.ctx,
			metav1.ListOptions{LabelSelector: labels, TimeoutSeconds: &timeout}); err != nil {
			return
		}
		if nsList, err = h.ListByLabel(labels); err != nil {
			return
		}
		if len(nsList) == 0 {
			isExist = false // namespace not exist
		} else {
			isExist = true // namespace exist
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
				log.Debug("watch namespace: bookmark.")
			case watch.Error:
				log.Debug("watch namespace: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch namespace: reconnect to kubernetes")
	}
}

// Watch watch namespace by name, alias to "WatchByName".
func (h *Handler) Watch(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
	return h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
}
