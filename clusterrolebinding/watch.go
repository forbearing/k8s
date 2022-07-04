package clusterrolebinding

import (
	log "github.com/sirupsen/logrus"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName  watch clusterrolebindings by name.
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
		if watcher, err = h.clientset.RbacV1().ClusterRoleBindings().Watch(h.ctx, listOptions); err != nil {
			return
		}
		if _, err = h.Get(name); err != nil {
			isExist = false // clusterrolebinding not exist
		} else {
			isExist = true // clusterrolebinding exist
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
				log.Debug("watch clusterrolebinding: bookmark.")
			case watch.Error:
				log.Debug("watch clusterrolebinding: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch clusterrolebinding: reconnect to kubernetes")
	}
}

// WatchByLabel watch clusterrolebindings by label.
func (h *Handler) WatchByLabel(labelSelector string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher                watch.Interface
		clusterrolebindingList *rbacv1.ClusterRoleBindingList
		timeout                = int64(0)
		isExist                bool
	)
	for {
		if watcher, err = h.clientset.RbacV1().ClusterRoleBindings().Watch(h.ctx,
			metav1.ListOptions{LabelSelector: labelSelector, TimeoutSeconds: &timeout}); err != nil {
			return
		}
		if clusterrolebindingList, err = h.List(labelSelector); err != nil {
			return
		}
		if len(clusterrolebindingList.Items) == 0 {
			isExist = false // clusterrolebinding not exist
		} else {
			isExist = true // clusterrolebinding exist
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
				log.Debug("watch clusterrolebinding: bookmark.")
			case watch.Error:
				log.Debug("watch clusterrolebinding: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch clusterrolebinding: reconnect to kubernetes")
	}
}

// Watch watch clusterrolebinding by name, alias to "WatchByName".
func (h *Handler) Watch(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) error {
	return h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
}
