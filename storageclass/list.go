package storageclass

import storagev1 "k8s.io/api/storage/v1"

// ListByLabel list storageclasses by labels.
func (h *Handler) ListByLabel(labels string) (*storagev1.StorageClassList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.StorageV1().StorageClasses().List(h.ctx, *listOptions)
}

// List list storageclasses by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*storagev1.StorageClassList, error) {
	return h.ListByLabel(labels)
}

// ListAll list all storageclasses in the k8s cluster.
func (h *Handler) ListAll() (*storagev1.StorageClassList, error) {
	return h.ListByLabel("")
}
