package daemonset

import (
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName watch daemonsets by name.
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
		if watcher, err = h.clientset.AppsV1().DaemonSets(h.namespace).Watch(h.ctx, listOptions); err != nil {
			return
		}
		if _, err = h.Get(name); err != nil {
			isExist = false // daemonsets not exist
		} else {
			isExist = true // daemonsets exist
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
				log.Debug("watch daemonset: bookmark")
			case watch.Error:
				log.Debug("watch daemonset: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch daemonset: reconnect to kubernetes")
		watcher.Stop()
	}
}

// WatchByLabel watch daemonsets by label.
func (h *Handler) WatchByLabel(labelSelector string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	var (
		watcher       watch.Interface
		daemonsetList *appsv1.DaemonSetList
		timeout       = int64(0)
		isExist       bool
	)
	for {
		if watcher, err = h.clientset.AppsV1().DaemonSets(h.namespace).Watch(h.ctx,
			metav1.ListOptions{LabelSelector: labelSelector, TimeoutSeconds: &timeout}); err != nil {
			return
		}
		if daemonsetList, err = h.List(labelSelector); err != nil {
			return
		}
		if len(daemonsetList.Items) == 0 {
			isExist = false // daemonset not exist
		} else {
			isExist = true // daemonset exist
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
				log.Debug("watch daemonset: bookmark")
			case watch.Error:
				log.Debug("watch daemonset: error")
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debug("watch daemonset: reconnect to kubernetes")
		watcher.Stop()
	}
}

// Watch watch daemonsets by name, alias to "WatchByName".
func (h *Handler) Watch(name string,
	addFunc, modifyFunc, deleteFunc func(x interface{}), x interface{}) (err error) {
	return h.WatchByName(name, addFunc, modifyFunc, deleteFunc, x)
}
