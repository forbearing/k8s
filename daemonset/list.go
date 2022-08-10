package daemonset

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all daemonsets in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*appsv1.DaemonSet, error) {
	return h.ListAll()
}

// ListByLabel list daemonsets by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*appsv1.DaemonSet, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	dsList, err := h.clientset.AppsV1().DaemonSets(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(dsList), nil
}

// ListByField list daemonsets by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*appsv1.DaemonSet, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	dsList, err := h.clientset.AppsV1().DaemonSets(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(dsList), nil
}

// ListByNamespace list all daemonsets in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*appsv1.DaemonSet, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all daemonsets in the k8s cluster.
func (h *Handler) ListAll() ([]*appsv1.DaemonSet, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(dsList *appsv1.DaemonSetList) []*appsv1.DaemonSet {
	var objList []*appsv1.DaemonSet
	for i := range dsList.Items {
		objList = append(objList, &dsList.Items[i])
	}
	return objList
}
