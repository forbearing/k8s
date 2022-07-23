package statefulset

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies statefulset from type string, []byte, *appsv1.StatefulSet,
// appsv1.StatefulSet, runtime.Object or map[string]interface{}.
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
	case runtime.Object:
		return h.ApplyFromObject(val)
	case map[string]interface{}:
		return h.ApplyFromUnstructured(val)
	default:
		return nil, ERR_TYPE_APPLY
	}
}

// ApplyFromFile applies statefulset from yaml file.
func (h *Handler) ApplyFromFile(filename string) (sts *appsv1.StatefulSet, err error) {
	sts, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if statefulset already exist, update it.
		sts, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply statefulset from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (sts *appsv1.StatefulSet, err error) {
	sts, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		sts, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies statefulset from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*appsv1.StatefulSet, error) {
	sts, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		return nil, fmt.Errorf("object is not *appsv1.StatefulSet")
	}
	return h.applyStatefulset(sts)
}

// ApplyFromUnstructured applies statefulset from map[string]interface{}.
func (h *Handler) ApplyFromUnstructured(u map[string]interface{}) (*appsv1.StatefulSet, error) {
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
