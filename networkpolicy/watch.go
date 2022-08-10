package networkpolicy

import (
	log "github.com/sirupsen/logrus"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName watch networkpolicyies by name.
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
		if watcher, err = h.clientset.NetworkingV1().NetworkPolicies(h.namespace).Watch(h.ctx, listOptions); err != nil {
			return
		}
		if _, err = h.Get(name); err != nil {
			isExist = false // networkpolicy not exist
		} else {
			isExist = true // networkpolicy exist
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
				log.Debug("watch networkpolicy: bookmark.")
			case watch.Error:
				log.Debug("watch networkpolicy: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch networkpolicy: reconnect to kubernetes")
	}
}

// WatchByLabel watch networkpolicies by labels.
func (h *Handler) WatchByLabel(labels string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher    watch.Interface
		netpolList []*networkingv1.NetworkPolicy
		timeout    = int64(0)
		isExist    bool
	)
	for {
		if watcher, err = h.clientset.NetworkingV1().NetworkPolicies(h.namespace).Watch(h.ctx,
			metav1.ListOptions{LabelSelector: labels, TimeoutSeconds: &timeout}); err != nil {
			return
		}
		if netpolList, err = h.ListByLabel(labels); err != nil {
			return
		}
		if len(netpolList) == 0 {
			isExist = false // networkpolicy not exist
		} else {
			isExist = true // networkpolicy exist
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
				log.Debug("watch networkpolicy: bookmark.")
			case watch.Error:
				log.Debug("watch networkpolicy: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch networkpolicy: reconnect to kubernetes")
	}
}

// Watch watch networkpolicies by name, alias to "WatchByName".
func (h *Handler) Watch(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	return h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
}
