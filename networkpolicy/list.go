package networkpolicy

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListByLabel list networkpolicies by labels.
func (h *Handler) ListByLabel(labels string) (*networkingv1.NetworkPolicyList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.NetworkingV1().NetworkPolicies(h.namespace).List(h.ctx, *listOptions)
}

// List list networkpolicies by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*networkingv1.NetworkPolicyList, error) {
	return h.ListByLabel(labels)
}

// ListByNamespace list networkpolicies by namespace.
func (h *Handler) ListByNamespace(namespace string) (*networkingv1.NetworkPolicyList, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all networkpolicies in the k8s cluster.
func (h *Handler) ListAll() (*networkingv1.NetworkPolicyList, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}
