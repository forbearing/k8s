package dynamic

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all k8s objects in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*unstructured.Unstructured, error) {
	return h.ListAll()
}

// ListByLabel list k8s objects by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*unstructured.Unstructured, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	if h.namespacedResource {
		unstructList, err := h.dynamicClient.Resource(h.gvr).Namespace(h.namespace).List(h.ctx, *listOptions)
		if err != nil {
			return nil, err
		}
		return extractList(unstructList), nil
	}
	unstructList, err := h.dynamicClient.Resource(h.gvr).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(unstructList), nil
}

// ListByField list k8s objects by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*unstructured.Unstructured, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	if h.namespacedResource {
		unstructList, err := h.dynamicClient.Resource(h.gvr).Namespace(h.namespace).List(h.ctx, *listOptions)
		if err != nil {
			return nil, err
		}
		return extractList(unstructList), nil
	}
	unstructList, err := h.dynamicClient.Resource(h.gvr).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(unstructList), nil
}

// ListByNamespace list all k8s objects in the specified namespace.
// It will return empty slice and error if this k8s object is cluster scope.
func (h *Handler) ListByNamespace(namespace string) ([]*unstructured.Unstructured, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	if h.namespacedResource {
		unstructList, err := h.dynamicClient.Resource(h.gvr).Namespace(namespace).List(h.ctx, *listOptions)
		if err != nil {
			return nil, err
		}
		return extractList(unstructList), nil
	}
	return nil, fmt.Errorf("%s is not namespace-scoped k8s resource", h.gvr)
}

// ListAll list all k8s objects in the k8s cluster.
func (h *Handler) ListAll() ([]*unstructured.Unstructured, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	if h.namespacedResource {
		unstructList, err := h.dynamicClient.Resource(h.gvr).Namespace(metav1.NamespaceAll).List(h.ctx, *listOptions)
		if err != nil {
			return nil, err
		}
		return extractList(unstructList), nil
	}
	unstructList, err := h.dynamicClient.Resource(h.gvr).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(unstructList), nil
}

// extractList
func extractList(unstructList *unstructured.UnstructuredList) []*unstructured.Unstructured {
	var objList []*unstructured.Unstructured
	for i := range unstructList.Items {
		objList = append(objList, &unstructList.Items[i])
	}
	return objList
}
