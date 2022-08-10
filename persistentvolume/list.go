package persistentvolume

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all persistentvolumes in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*corev1.PersistentVolume, error) {
	return h.ListAll()
}

// ListByLabel list persistentvolumes by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*corev1.PersistentVolume, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	pvList, err := h.clientset.CoreV1().PersistentVolumes().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(pvList), nil
}

// ListByField list persistentvolumes by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*corev1.PersistentVolume, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	pvList, err := h.clientset.CoreV1().PersistentVolumes().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(pvList), nil
}

// ListAll list all persistentvolumes in the k8s cluster.
func (h *Handler) ListAll() ([]*corev1.PersistentVolume, error) {
	return h.ListByLabel("")
}

// extractList
func extractList(pvList *corev1.PersistentVolumeList) []*corev1.PersistentVolume {
	var objList []*corev1.PersistentVolume
	for i := range pvList.Items {
		objList = append(objList, &pvList.Items[i])
	}
	return objList
}
