package persistentvolume

import corev1 "k8s.io/api/core/v1"

// ListByLabel list persistentvolumes by labels.
func (h *Handler) ListByLabel(labels string) (*corev1.PersistentVolumeList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.CoreV1().PersistentVolumes().List(h.ctx, *listOptions)
}

// List list persistentvolumes by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*corev1.PersistentVolumeList, error) {
	return h.ListByLabel(labels)
}

// ListAll list all persistentvolumes in the k8s cluster.
func (h *Handler) ListAll() (*corev1.PersistentVolumeList, error) {
	return h.ListByLabel("")
}
