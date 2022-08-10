package persistentvolumeclaim

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all persistentvolumeclaims in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*corev1.PersistentVolumeClaim, error) {
	return h.ListAll()
}

// ListByLabel list persistentvolumeclaims by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*corev1.PersistentVolumeClaim, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	pvcList, err := h.clientset.CoreV1().PersistentVolumeClaims(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(pvcList), nil
}

// ListByField list persistentvolumeclaims by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*corev1.PersistentVolumeClaim, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	pvcList, err := h.clientset.CoreV1().PersistentVolumeClaims(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(pvcList), nil
}

// ListByNamespace list all persistentvolumeclaims in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*corev1.PersistentVolumeClaim, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all persistentvolumeclaims in the k8s cluster.
func (h *Handler) ListAll() ([]*corev1.PersistentVolumeClaim, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(pvcList *corev1.PersistentVolumeClaimList) []*corev1.PersistentVolumeClaim {
	var objList []*corev1.PersistentVolumeClaim
	for i := range pvcList.Items {
		objList = append(objList, &pvcList.Items[i])
	}
	return objList
}
