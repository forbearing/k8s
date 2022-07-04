package serviceaccount

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListByLabel list serviceaccounts by labels.
func (h *Handler) ListByLabel(labels string) (*corev1.ServiceAccountList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.CoreV1().ServiceAccounts(h.namespace).List(h.ctx, *listOptions)
}

// List list serviceaccounts by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*corev1.ServiceAccountList, error) {
	return h.ListByLabel(labels)
}

// ListByNamespace list serviceaccounts by namespace.
func (h *Handler) ListByNamespace(namespace string) (*corev1.ServiceAccountList, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all serviceaccounts in the k8s cluster.
func (h *Handler) ListAll(namespace string) (*corev1.ServiceAccountList, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}
