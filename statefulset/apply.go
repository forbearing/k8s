package statefulset

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies statefulset from type string, []byte, *appsv1.StatefulSet,
// appsv1.StatefulSet, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*appsv1.StatefulSet, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *appsv1.StatefulSet:
		return h.ApplyFromObject(val)
	case appsv1.StatefulSet:
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

// ApplyFromFile applies statefulset from yaml or json file.
func (h *Handler) ApplyFromFile(filename string) (sts *appsv1.StatefulSet, err error) {
	sts, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if statefulset already exist, update it.
		sts, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply statefulset from bytes data.
func (h *Handler) ApplyFromBytes(data []byte) (sts *appsv1.StatefulSet, err error) {
	sts, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		sts, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies statefulset from metav1.Object or runtime.Object.
func (h *Handler) ApplyFromObject(obj interface{}) (*appsv1.StatefulSet, error) {
	sts, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.StatefulSet")
	}
	return h.applyStatefulset(sts)
}

// ApplyFromUnstructured applies statefulset from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sts)
	if err != nil {
		return nil, err
	}
	return h.applyStatefulset(sts)
}

// ApplyFromMap applies statefulset from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sts)
	if err != nil {
		return nil, err
	}
	return h.applyStatefulset(sts)
}

// applyStatefulset
func (h *Handler) applyStatefulset(sts *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	_, err := h.createStatefulset(sts)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateStatefulset(sts)
	}
	return sts, err
}
