package replicaset

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies replicaset from type string, []byte, *appsv1.ReplicaSet,
// appsv1.ReplicaSet, runtime.Object or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*appsv1.ReplicaSet, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *appsv1.ReplicaSet:
		return h.ApplyFromObject(val)
	case appsv1.ReplicaSet:
		return h.ApplyFromObject(&val)
	case runtime.Object:
		return h.ApplyFromObject(val)
	case map[string]interface{}:
		return h.ApplyFromUnstructured(val)
	default:
		return nil, ERR_TYPE_APPLY
	}
}

// ApplyFromFile applies replicaset from yaml file.
func (h *Handler) ApplyFromFile(filename string) (rs *appsv1.ReplicaSet, err error) {
	rs, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if replicaset already exist, update it.
		rs, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply replicaset from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (rs *appsv1.ReplicaSet, err error) {
	rs, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		rs, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies replicaset from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*appsv1.ReplicaSet, error) {
	rs, ok := obj.(*appsv1.ReplicaSet)
	if !ok {
		return nil, fmt.Errorf("object is not *appsv1.ReplicaSet")
	}
	return h.applyReplicaset(rs)
}

// ApplyFromUnstructured applies replicaset from map[string]interface{}.
func (h *Handler) ApplyFromUnstructured(u map[string]interface{}) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rs)
	if err != nil {
		return nil, err
	}
	return h.applyReplicaset(rs)
}

// applyReplicaset
func (h *Handler) applyReplicaset(rs *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error) {
	_, err := h.createReplicaset(rs)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateReplicaset(rs)
	}
	return rs, err
}
