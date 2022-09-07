package dynamic

import (
	utilrestmapper "github.com/forbearing/k8s/util/restmapper"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

// WatchByName watch a single k8s resource
// You should always specify the GroupVersionKind with WithGVK() method.
//
// Object as the parameter of addFunc, modifyFunc, deleteFunc:
//  * If Event.Type is Added or Modified: the new state of the object.
//  * If Event.Type is Deleted: the state of the object immediately before deletion.
//  * If Event.Type is Bookmark: the object (instance of a type being watched) where
//    only ResourceVersion field is set. On successful restart of watch from a
//    bookmark resourceVersion, client is guaranteed to not get repeat event
//    nor miss any events.
//  * If Event.Type is Error: *api.Status is recommended; other types may make sense
//    depending on context.
func (h *Handler) WatchByName(name string, addFunc, modifyFunc, deleteFunc func(obj interface{})) error {
	var (
		err          error
		watcher      watch.Interface
		gvr          schema.GroupVersionResource
		isNamespaced bool
	)

	listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
	listOptions.TimeoutSeconds = new(int64)
	if gvr, err = utilrestmapper.GVKToGVR(h.restMapper, h.gvk); err != nil {
		return err
	}
	if isNamespaced, err = utilrestmapper.IsNamespaced(h.restMapper, h.gvk); err != nil {
		return err
	}
	for {
		if isNamespaced {
			if watcher, err = h.dynamicClient.Resource(gvr).Namespace(h.namespace).Watch(h.ctx, listOptions); err != nil {
				return err
			}
		}
		if watcher, err = h.dynamicClient.Resource(gvr).Watch(h.ctx, listOptions); err != nil {
			return err
		}
		for event := range watcher.ResultChan() {
			switch event.Type {
			case watch.Added:
				addFunc(event.Object)
			case watch.Modified:
				modifyFunc(event.Object)
			case watch.Deleted:
				deleteFunc(event.Object)
			case watch.Bookmark:
				log.Debugf("watch %s: bookmark.", gvr.String())
			case watch.Error:
				log.Debugf("watch %s: error", gvr.String())
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debugf("watch %s: reconnect to kubernetes", gvr.String())
		watcher.Stop()
	}
}
