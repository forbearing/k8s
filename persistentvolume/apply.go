package persistentvolume

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies persistentvolume from type string, []byte, *corev1.PersistentVolume,
// corev1.PersistentVolume, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*corev1.PersistentVolume, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *corev1.PersistentVolume:
		return h.ApplyFromObject(val)
	case corev1.PersistentVolume:
		return h.ApplyFromObject(&val)
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	case runtime.Object:
		return h.ApplyFromObject(val)
	default:
		return nil, ErrInvalidApplyType
	}
}

// ApplyFromFile applies persistentvolume from yaml file.
func (h *Handler) ApplyFromFile(filename string) (pv *corev1.PersistentVolume, err error) {
	pv, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if persistentvolume already exist, update it.
		pv, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply persistentvolume from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (pv *corev1.PersistentVolume, err error) {
	pv, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		pv, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies persistentvolume from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*corev1.PersistentVolume, error) {
	pv, ok := obj.(*corev1.PersistentVolume)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.PersistentVolume")
	}
	return h.applyPV(pv)
}

// ApplyFromUnstructured applies persistentvolume from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*corev1.PersistentVolume, error) {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pv)
	if err != nil {
		return nil, err
	}
	return h.applyPV(pv)
}

// ApplyFromMap applies persistentvolume from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*corev1.PersistentVolume, error) {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pv)
	if err != nil {
		return nil, err
	}
	return h.applyPV(pv)
}

// applyPV
func (h *Handler) applyPV(pv *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	_, err := h.createPV(pv)
	if k8serrors.IsAlreadyExists(err) {
		return h.updatePV(pv)
	}
	return pv, err
}
