package ingressclass

import (
	log "github.com/sirupsen/logrus"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName watch ingressclass by name.
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
		if watcher, err = h.clientset.NetworkingV1().IngressClasses().Watch(h.ctx, listOptions); err != nil {
			return
		}
		if _, err = h.Get(name); err != nil {
			isExist = false // configmap not exist
		} else {
			isExist = true // configmap exist
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
				log.Debug("watch ingressclass: bookmark.")
			case watch.Error:
				log.Debug("watch ingressclass: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch ingressclass: reconnect to kubernetes")
	}
}

// WatchByLabel watch ingressclass by labels.
func (h *Handler) WatchByLabel(labels string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher  watch.Interface
		ingcList []*networkingv1.IngressClass
		timeout  = int64(0)
		isExist  bool
	)
	for {
		if watcher, err = h.clientset.NetworkingV1().IngressClasses().Watch(h.ctx,
			metav1.ListOptions{LabelSelector: labels, TimeoutSeconds: &timeout}); err != nil {
			return
		}
		if ingcList, err = h.ListByLabel(labels); err != nil {
			return
		}
		if len(ingcList) == 0 {
			isExist = false // ingressclass not exist
		} else {
			isExist = true // ingressclass exist
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
				log.Debug("watch ingressclass: bookmark.")
			case watch.Error:
				log.Debug("watch ingressclass: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch ingressclass: reconnect to kubernetes")
	}
}

// Watch watch ingressclass by name, alias to "WatchByName".
func (h *Handler) Watch(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	return h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
}
