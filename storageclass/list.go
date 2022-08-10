package storageclass

import (
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all storageclasses in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*storagev1.StorageClass, error) {
	return h.ListAll()
}

// ListByLabel list storageclasses by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*storagev1.StorageClass, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	scList, err := h.clientset.StorageV1().StorageClasses().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(scList), nil
}

// ListByField list storageclasses by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*storagev1.StorageClass, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	scList, err := h.clientset.StorageV1().StorageClasses().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(scList), nil
}

// ListAll list all storageclasses in the k8s cluster.
func (h *Handler) ListAll() ([]*storagev1.StorageClass, error) {
	return h.ListByLabel("")
}

// extractList
func extractList(scList *storagev1.StorageClassList) []*storagev1.StorageClass {
	var objList []*storagev1.StorageClass
	for i := range scList.Items {
		objList = append(objList, &scList.Items[i])
	}
	return objList
}
