package statefulset

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all statefulsets in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*appsv1.StatefulSet, error) {
	return h.ListAll()
}

// ListByLabel list statefulsets by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*appsv1.StatefulSet, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	stsList, err := h.clientset.AppsV1().StatefulSets(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(stsList), nil
}

// ListByField list statefulsets by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*appsv1.StatefulSet, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	stsList, err := h.clientset.AppsV1().StatefulSets(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(stsList), nil
}

// ListByNamespace list all statefulsets in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*appsv1.StatefulSet, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all statefulsets in the k8s cluster.
func (h *Handler) ListAll() ([]*appsv1.StatefulSet, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(stsList *appsv1.StatefulSetList) []*appsv1.StatefulSet {
	var objList []*appsv1.StatefulSet
	for i := range stsList.Items {
		objList = append(objList, &stsList.Items[i])
	}
	return objList
}
