package storageclass

import (
	log "github.com/sirupsen/logrus"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName watch storageclass by name.
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
		if watcher, err = h.clientset.StorageV1().StorageClasses().Watch(h.ctx, listOptions); err != nil {
			return
		}
		if _, err = h.Get(name); err != nil {
			isExist = false // sc not exist
		} else {
			isExist = true // sc exist
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
				log.Debug("watch storageclass: bookmark.")
			case watch.Error:
				log.Debug("watch storageclass: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch storageclass: reconnect to kubernetes")
	}
}

// WatchByLabel watch storageclass by labels.
func (h *Handler) WatchByLabel(labels string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher watch.Interface
		scList  []*storagev1.StorageClass
		timeout = int64(0)
		isExist bool
	)
	for {
		if watcher, err = h.clientset.StorageV1().StorageClasses().Watch(h.ctx,
			metav1.ListOptions{LabelSelector: labels, TimeoutSeconds: &timeout}); err != nil {
			return
		}
		if scList, err = h.ListByLabel(labels); err != nil {
			return
		}
		if len(scList) == 0 {
			isExist = false // sc not exist
		} else {
			isExist = true // sc exist
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
				log.Debug("watch storageclass: bookmark.")
			case watch.Error:
				log.Debug("watch storageclass: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch storageclass: reconnect to kubernetes")
	}
}

// Watch watch storageclass by name, alias to "WatchByName".
func (h *Handler) Watch(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	return h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
}
