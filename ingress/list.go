package ingress

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListByLabel list ingresses by labels.
func (h *Handler) ListByLabel(labels string) (*networkingv1.IngressList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.NetworkingV1().Ingresses(h.namespace).List(h.ctx, *listOptions)
}

// List list ingresses by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*networkingv1.IngressList, error) {
	return h.ListByLabel(labels)
}

// ListByNamespace list ingresses by namespace.
func (h *Handler) ListByNamespace(namespace string) (*networkingv1.IngressList, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all ingresses in the k8s cluster.
func (h *Handler) ListAll() (*networkingv1.IngressList, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}
