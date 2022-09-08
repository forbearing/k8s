package dynamic

import (
	utilrestmapper "github.com/forbearing/k8s/util/restmapper"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
)

// Watch watch all k8s resource with the specified kind.
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
func (h *Handler) Watch(addFunc, modifyFunc, deleteFunc func(obj interface{})) error {
	return h.WithNamespace(metav1.NamespaceAll).WatchByLabel("", addFunc, modifyFunc, deleteFunc)
}

// WatchByNamespace watch all k8s resource with the specified kind in the specified namespace.
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
func (h *Handler) WatchByNamespace(namespace string, addFunc, modifyFunc, deleteFunc func(obj interface{})) error {
	if len(namespace) == 0 {
		namespace = metav1.NamespaceDefault
	}
	return h.WithNamespace(namespace).WatchByLabel("", addFunc, modifyFunc, deleteFunc)
}

// WatchByName watch a single k8s resource with the specified Kind.
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
	listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
	listOptions.TimeoutSeconds = new(int64)
	return h.watchUnstructuredObj(listOptions, addFunc, modifyFunc, deleteFunc)
}

// WatchByLabel watch a single or multiple k8s resource with the specified Kind
// and selected by the labels. Multiple labels are separated by ",",
// label key and value conjunctaed by "=".
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
func (h *Handler) WatchByLabel(labels string, addFunc, modifyFunc, deleteFunc func(obj interface{})) error {
	return h.watchUnstructuredObj(
		metav1.ListOptions{LabelSelector: labels, TimeoutSeconds: new(int64)},
		addFunc, modifyFunc, deleteFunc)
}

// WatchByField watch a single or multiple k8s resources with specified Kind
// and selected by the field.
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
func (h *Handler) WatchByField(field string, addFunc, modifyFunc, deleteFunc func(obj interface{})) error {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return err
	}
	listOptions := metav1.ListOptions{FieldSelector: fieldSelector.String(), TimeoutSeconds: new(int64)}
	return h.watchUnstructuredObj(listOptions, addFunc, modifyFunc, deleteFunc)
}

// watchUnstructuredObj watch k8s object according to the listOptions.
func (h *Handler) watchUnstructuredObj(listOptions metav1.ListOptions,
	addFunc, modifyFunc, deleteFunc func(obj interface{})) (err error) {

	var gvr schema.GroupVersionResource
	var isNamespaced bool
	if gvr, err = utilrestmapper.GVKToGVR(h.restMapper, h.gvk); err != nil {
		return err
	}
	if isNamespaced, err = utilrestmapper.IsNamespaced(h.restMapper, h.gvk); err != nil {
		return err
	}

	var watcher watch.Interface
	// If event channel is closed, it means the server has closed the connection,
	// reconnect to kubernetes API server.
	for {
		if isNamespaced {
			if watcher, err = h.dynamicClient.Resource(gvr).Namespace(h.namespace).Watch(h.ctx, listOptions); err != nil {
				return err
			}
		} else {
			if watcher, err = h.dynamicClient.Resource(gvr).Watch(h.ctx, listOptions); err != nil {
				return err
			}
		}
		// Kubernetes retains the resource event history, which includes this
		// initial event, so that when our program first start, we are automatically
		// notified of the deployment existence and current state.
		// There we will not ignore the first resource added event.
		for event := range watcher.ResultChan() {
			switch event.Type {
			case watch.Added:
				addFunc(event.Object)
			case watch.Modified:
				modifyFunc(event.Object)
			case watch.Deleted:
				deleteFunc(event.Object)
			case watch.Bookmark:
				log.Debugf("watch %s: bookmark", gvr.Resource)
			case watch.Error:
				log.Debugf("watch %s: error", gvr.Resource)
			}
		}
		// If event channel is closed, it means the server has closed the connection
		log.Debugf("watch %s: reconnect to kubernetes", gvr.Resource)
		watcher.Stop()
	}
}
