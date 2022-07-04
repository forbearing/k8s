package ingressclass

import networkingv1 "k8s.io/api/networking/v1"

// ListByLabel list ingressclasses by labels.
func (h *Handler) ListByLabel(labels string) (*networkingv1.IngressClassList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.NetworkingV1().IngressClasses().List(h.ctx, *listOptions)
}

// List list ingressclasses by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*networkingv1.IngressClassList, error) {
	return h.ListByLabel(labels)
}

// ListAll list all ingressclasses in the k8s cluster.
func (h *Handler) ListAll(labels string) (*networkingv1.IngressClassList, error) {
	return h.ListByLabel("")
}
