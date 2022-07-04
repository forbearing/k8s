package clusterrole

import (
	log "github.com/sirupsen/logrus"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName watch clusterroles by name.
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
		if watcher, err = h.clientset.RbacV1().ClusterRoles().Watch(h.ctx, listOptions); err != nil {
			return
		}
		if _, err = h.Get(name); err != nil {
			isExist = false // clusterroles not exist
		} else {
			isExist = true // clusterroles exist
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
				log.Debug("watch clusterrole: bookmark.")
			case watch.Error:
				log.Debug("watch clusterrole: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch clusterrole: reconnect to kubernetes")
	}
}

// WatchByLabel watch clusterroles by label.
func (h *Handler) WatchByLabel(labelSelector string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher         watch.Interface
		clusterroleList *rbacv1.ClusterRoleList
		timeout         = int64(0)
		isExist         bool
	)
	for {
		if watcher, err = h.clientset.RbacV1().ClusterRoles().Watch(h.ctx,
			metav1.ListOptions{LabelSelector: labelSelector, TimeoutSeconds: &timeout}); err != nil {
			return
		}
		if clusterroleList, err = h.List(labelSelector); err != nil {
			return
		}
		if len(clusterroleList.Items) == 0 {
			isExist = false // clusterrole not exist
		} else {
			isExist = true // clusterrole exist
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
				log.Debug("watch clusterrole: bookmark.")
			case watch.Error:
				log.Debug("watch clusterrole: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch clusterrole: reconnect to kubernetes")
	}
}

// Watch watch clusterroles by name, alias to "WatchByName".
func (h *Handler) Watch(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	return h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
}
