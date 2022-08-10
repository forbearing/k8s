package serviceaccount

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all serviceaccounts in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*corev1.ServiceAccount, error) {
	return h.ListAll()
}

// ListByLabel list serviceaccounts by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*corev1.ServiceAccount, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	saList, err := h.clientset.CoreV1().ServiceAccounts(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(saList), nil
}

// ListByField list serviceaccounts by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*corev1.ServiceAccount, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	saList, err := h.clientset.CoreV1().ServiceAccounts(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(saList), nil
}

// ListByNamespace list all serviceaccounts in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*corev1.ServiceAccount, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all serviceaccounts in the k8s cluster.
func (h *Handler) ListAll() ([]*corev1.ServiceAccount, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(saList *corev1.ServiceAccountList) []*corev1.ServiceAccount {
	var objList []*corev1.ServiceAccount
	for i := range saList.Items {
		objList = append(objList, &saList.Items[i])
	}
	return objList
}
