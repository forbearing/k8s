package node

import corev1 "k8s.io/api/core/v1"

// GetByName get node by name.
func (h *Handler) GetByName(name string) (*corev1.Node, error) {
	return h.clientset.CoreV1().Nodes().Get(h.ctx, name, h.Options.GetOptions)
}

// Get get node by name, alias to "GetByName".
func (h *Handler) Get(name string) (*corev1.Node, error) {
	return h.GetByName(name)
}
