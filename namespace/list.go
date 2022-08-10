package namespace

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all namespaces in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*corev1.Namespace, error) {
	return h.ListAll()
}

// ListByLabel list namespaces by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*corev1.Namespace, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	nsList, err := h.clientset.CoreV1().Namespaces().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(nsList), nil
}

// ListByField list namespaces by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*corev1.Namespace, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	nsList, err := h.clientset.CoreV1().Namespaces().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(nsList), nil
}

// ListAll list all namespaces in the k8s cluster.
func (h *Handler) ListAll() ([]*corev1.Namespace, error) {
	return h.ListByLabel("")
}

// extractList
func extractList(nsList *corev1.NamespaceList) []*corev1.Namespace {
	var objList []*corev1.Namespace
	for i := range nsList.Items {
		objList = append(objList, &nsList.Items[i])
	}
	return objList
}
