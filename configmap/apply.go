package configmap

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies configmap from type string, []byte, *corev1.ConfigMap,
// corev1.ConfigMap, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*corev1.ConfigMap, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *corev1.ConfigMap:
		return h.ApplyFromObject(val)
	case corev1.ConfigMap:
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

// ApplyFromFile applies configmap from yaml or json file.
func (h *Handler) ApplyFromFile(filename string) (cm *corev1.ConfigMap, err error) {
	cm, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if configmap already exist, update it.
		cm, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply configmap from bytes data.
func (h *Handler) ApplyFromBytes(data []byte) (cm *corev1.ConfigMap, err error) {
	cm, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		cm, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies configmap from metav1.Object or runtime.Object.
func (h *Handler) ApplyFromObject(obj interface{}) (*corev1.ConfigMap, error) {
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.ConfigMap")
	}
	return h.applyConfigmap(cm)
}

// ApplyFromUnstructured applies configmap from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cm)
	if err != nil {
		return nil, err
	}
	return h.applyConfigmap(cm)
}

// ApplyFromMap applies configmap from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cm)
	if err != nil {
		return nil, err
	}
	return h.applyConfigmap(cm)
}

// applyConfigmap
func (h *Handler) applyConfigmap(cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	_, err := h.createConfigmap(cm)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateConfigmap(cm)
	}
	return cm, err
}
