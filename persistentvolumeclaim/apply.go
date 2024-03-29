package persistentvolumeclaim

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies persistentvolumeclaim from type string, []byte, *corev1.PersistentVolumeClaim,
// corev1.PersistentVolumeClaim, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*corev1.PersistentVolumeClaim, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *corev1.PersistentVolumeClaim:
		return h.ApplyFromObject(val)
	case corev1.PersistentVolumeClaim:
		return h.ApplyFromObject(&val)
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	case metav1.Object, runtime.Object:
		return h.ApplyFromObject(val)
	default:
		return nil, ErrInvalidApplyType
	}
}

// ApplyFromFile applies persistentvolumeclaim from yaml or json file.
func (h *Handler) ApplyFromFile(filename string) (pvc *corev1.PersistentVolumeClaim, err error) {
	pvc, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if persistentvolumeclaim already exist, update it.
		pvc, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply persistentvolumeclaim from bytes data.
func (h *Handler) ApplyFromBytes(data []byte) (pvc *corev1.PersistentVolumeClaim, err error) {
	pvc, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		pvc, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies persistentvolumeclaim from metav1.Object or runtime.Object.
func (h *Handler) ApplyFromObject(obj interface{}) (*corev1.PersistentVolumeClaim, error) {
	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.PersistentVolumeClaim")
	}
	return h.applyPVC(pvc)
}

// ApplyFromUnstructured applies persistentvolumeclaim from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pvc)
	if err != nil {
		return nil, err
	}
	return h.applyPVC(pvc)
}

// ApplyFromMap applies persistentvolumeclaim from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pvc)
	if err != nil {
		return nil, err
	}
	return h.applyPVC(pvc)
}

// applyPVC
func (h *Handler) applyPVC(pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	_, err := h.createPVC(pvc)
	if k8serrors.IsAlreadyExists(err) {
		return h.updatePVC(pvc)
	}
	return pvc, err
}
