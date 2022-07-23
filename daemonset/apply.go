package daemonset

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies daemonset from type string, []byte, *appsv1.DaemonSet,
// appsv1.DaemonSet, runtime.Object or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*appsv1.DaemonSet, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *appsv1.DaemonSet:
		return h.ApplyFromObject(val)
	case appsv1.DaemonSet:
		return h.ApplyFromObject(&val)
	case runtime.Object:
		return h.ApplyFromObject(val)
	case map[string]interface{}:
		return h.ApplyFromUnstructured(val)
	default:
		return nil, ERR_TYPE_APPLY
	}
}

// ApplyFromFile applies daemonset from yaml file.
func (h *Handler) ApplyFromFile(filename string) (ds *appsv1.DaemonSet, err error) {
	ds, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if daemonset already exist, update it.
		ds, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply daemonset from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (ds *appsv1.DaemonSet, err error) {
	ds, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		ds, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies daemonset from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*appsv1.DaemonSet, error) {
	ds, ok := obj.(*appsv1.DaemonSet)
	if !ok {
		return nil, fmt.Errorf("object is not *appsv1.DaemonSet")
	}
	return h.applyDaemonset(ds)
}

// ApplyFromUnstructured applies daemonset from map[string]interface{}.
func (h *Handler) ApplyFromUnstructured(u map[string]interface{}) (*appsv1.DaemonSet, error) {
	ds := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ds)
	if err != nil {
		return nil, err
	}
	return h.applyDaemonset(ds)
}

// applyDaemonset
func (h *Handler) applyDaemonset(ds *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	_, err := h.createDaemonset(ds)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateDaemonset(ds)
	}
	return ds, err
}
