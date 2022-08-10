package replicaset

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all replicasets in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*appsv1.ReplicaSet, error) {
	return h.ListAll()
}

// ListByLabel list replicasets by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*appsv1.ReplicaSet, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	rsList, err := h.clientset.AppsV1().ReplicaSets(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(rsList), nil
}

// ListByField list replicasets by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*appsv1.ReplicaSet, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	rsList, err := h.clientset.AppsV1().ReplicaSets(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(rsList), nil
}

// ListByNamespace list all replicasets in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*appsv1.ReplicaSet, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all replicasets in the k8s cluster.
func (h *Handler) ListAll() ([]*appsv1.ReplicaSet, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(rsList *appsv1.ReplicaSetList) []*appsv1.ReplicaSet {
	var objList []*appsv1.ReplicaSet
	for i := range rsList.Items {
		objList = append(objList, &rsList.Items[i])
	}
	return objList
}
