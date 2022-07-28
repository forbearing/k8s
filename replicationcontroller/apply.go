package replicationcontroller

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies replicationcontroller from type string, []byte,
// *corev1.ReplicationController, corev1.ReplicationController, runtime.Object,
// *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*corev1.ReplicationController, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *corev1.ReplicationController:
		return h.ApplyFromObject(val)
	case corev1.ReplicationController:
		return h.ApplyFromObject(&val)
	case runtime.Object:
		return h.ApplyFromObject(val)
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	default:
		return nil, ERR_TYPE_APPLY
	}
}

// ApplyFromFile applies replicationcontroller from yaml file.
func (h *Handler) ApplyFromFile(filename string) (rc *corev1.ReplicationController, err error) {
	rc, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if replicationcontroller already exist, update it.
		rc, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply replicationcontroller from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (rc *corev1.ReplicationController, err error) {
	rc, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		rc, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies replicationcontroller from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*corev1.ReplicationController, error) {
	rc, ok := obj.(*corev1.ReplicationController)
	if !ok {
		return nil, fmt.Errorf("object is not *corev1.ReplicationController")
	}
	return h.applyRS(rc)
}

// ApplyFromUnstructured applies replicationcontroller from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*corev1.ReplicationController, error) {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rc)
	if err != nil {
		return nil, err
	}
	return h.applyRS(rc)
}

// ApplyFromMap applies replicationcontroller from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*corev1.ReplicationController, error) {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rc)
	if err != nil {
		return nil, err
	}
	return h.applyRS(rc)
}

// applyRS
func (h *Handler) applyRS(rc *corev1.ReplicationController) (*corev1.ReplicationController, error) {
	_, err := h.createRS(rc)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateRS(rc)
	}
	return rc, err
}
