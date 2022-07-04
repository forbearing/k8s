package service

import (
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName watch services by name.
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
		if watcher, err = h.clientset.CoreV1().Services(h.namespace).Watch(h.ctx, listOptions); err != nil {
			logrus.Error(err)
			return
		}
		if _, err = h.Get(name); err != nil {
			isExist = false // service not exist
		} else {
			isExist = true // service exist
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
				log.Debug("watch service: bookmark.")
			case watch.Error:
				log.Debug("watch service: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch service: reconnect to kubernetes")
	}
}

// WatchByLabel watch services by label.
func (h *Handler) WatchByLabel(labelSelector string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher     watch.Interface
		serviceList *corev1.ServiceList
		timeout     = int64(0)
		isExist     bool
	)
	for {
		if watcher, err = h.clientset.CoreV1().Services(h.namespace).Watch(h.ctx,
			metav1.ListOptions{LabelSelector: labelSelector, TimeoutSeconds: &timeout}); err != nil {
			logrus.Error(err)
			return
		}
		if serviceList, err = h.List(labelSelector); err != nil {
			logrus.Error(err)
			return
		}
		if len(serviceList.Items) == 0 {
			isExist = false // service not exist
		} else {
			isExist = true // service exist
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
				log.Debug("watch service: bookmark.")
			case watch.Error:
				log.Debug("watch service: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch service: reconnect to kubernetes")
	}
}

// Watch watch services by name, alias to "WatchByName".
func (h *Handler) Watch(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	return h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
}
