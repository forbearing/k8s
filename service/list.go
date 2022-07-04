package service

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListByLabel list services by labels.
func (h *Handler) ListByLabel(labels string) (*corev1.ServiceList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.CoreV1().Services(h.namespace).List(h.ctx, *listOptions)
}

// List list services by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*corev1.ServiceList, error) {
	return h.ListByLabel(labels)
}

// ListByNamespace list services by labels, alias to "ListByLabel".
func (h *Handler) ListByNamespace(namespace string) (*corev1.ServiceList, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all services in the k8s cluster.
func (h *Handler) ListAll() (*corev1.ServiceList, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}
