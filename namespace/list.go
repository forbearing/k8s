package namespace

import corev1 "k8s.io/api/core/v1"

// ListByLabel list namespaces by labels.
func (h *Handler) ListByLabel(labels string) (*corev1.NamespaceList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.CoreV1().Namespaces().List(h.ctx, *listOptions)
}

// List list namespaces by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*corev1.NamespaceList, error) {
	return h.ListByLabel(labels)
}

// ListAll list all namespaces in the k8s cluster.
func (h *Handler) ListAll(labels string) (*corev1.NamespaceList, error) {
	return h.ListByLabel("")
}
