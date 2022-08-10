package configmap

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all configmaps in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*corev1.ConfigMap, error) {
	return h.ListAll()
}

// ListByLabel list configmaps by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*corev1.ConfigMap, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	cmList, err := h.clientset.CoreV1().ConfigMaps(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(cmList), nil
}

// ListByField list configmaps by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*corev1.ConfigMap, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	cmList, err := h.clientset.CoreV1().ConfigMaps(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(cmList), nil
}

// ListByNamespace list all configmaps in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*corev1.ConfigMap, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all configmaps in the k8s cluster.
func (h *Handler) ListAll() ([]*corev1.ConfigMap, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(cmList *corev1.ConfigMapList) []*corev1.ConfigMap {
	var objList []*corev1.ConfigMap
	for i := range cmList.Items {
		objList = append(objList, &cmList.Items[i])
	}
	return objList
}
