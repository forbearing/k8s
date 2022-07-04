package node

import corev1 "k8s.io/api/core/v1"

// ListByLabel list nodes by labels.
func (h *Handler) ListByLabel(labels string) (*corev1.NodeList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.CoreV1().Nodes().List(h.ctx, *listOptions)
}

// List list nodes by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*corev1.NodeList, error) {
	return h.ListByLabel(labels)
}

// ListAll list all nodes.
func (h *Handler) ListAll() (*corev1.NodeList, error) {
	return h.ListByLabel("")
}
