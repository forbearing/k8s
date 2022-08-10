package rolebinding

import (
	log "github.com/sirupsen/logrus"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName watch rolebindings by name.
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
		if watcher, err = h.clientset.RbacV1().RoleBindings(h.namespace).Watch(h.ctx, listOptions); err != nil {
			return
		}
		if _, err = h.Get(name); err != nil {
			isExist = false // role not exist
		} else {
			isExist = true // role exist
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
				log.Debug("watch rolebinding: bookmark.")
			case watch.Error:
				log.Debug("watch rolebinding: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch rolebinding: reconnect to kubernetes")
	}
}

// WatchByLabel watch rolebindings by labels.
func (h *Handler) WatchByLabel(labels string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher watch.Interface
		rbList  []*rbacv1.RoleBinding
		timeout = int64(0)
		isExist bool
	)
	for {
		if watcher, err = h.clientset.RbacV1().RoleBindings(h.namespace).Watch(h.ctx,
			metav1.ListOptions{LabelSelector: labels, TimeoutSeconds: &timeout}); err != nil {
			return
		}
		if rbList, err = h.ListByLabel(labels); err != nil {
			return
		}
		if len(rbList) == 0 {
			isExist = false // rolebinding not exist
		} else {
			isExist = true // rolebinding exist
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
				log.Debug("watch rolebinding: bookmark.")
			case watch.Error:
				log.Debug("watch rolebinding: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch rolebinding: reconnect to kubernetes")
	}
}

// Watch watch rolebindings by name, alias to "WatchByName".
func (h *Handler) Watch(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	return h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
}
