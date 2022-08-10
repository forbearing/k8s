package service

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all services in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*corev1.Service, error) {
	return h.ListAll()
}

// ListByLabel list services by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*corev1.Service, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	svcList, err := h.clientset.CoreV1().Services(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(svcList), nil
}

// ListByField list services by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*corev1.Service, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	svcList, err := h.clientset.CoreV1().Services(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(svcList), nil
}

// ListByNamespace list all services in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*corev1.Service, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all services in the k8s cluster.
func (h *Handler) ListAll() ([]*corev1.Service, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(svcList *corev1.ServiceList) []*corev1.Service {
	var objList []*corev1.Service
	for i := range svcList.Items {
		objList = append(objList, &svcList.Items[i])
	}
	return objList
}
