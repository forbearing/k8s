package clusterrole

import (
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies clusterrole from type string, []byte, *rbacv1.ClusterRole,
// rbacv1.ClusterRole, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*rbacv1.ClusterRole, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *rbacv1.ClusterRole:
		return h.ApplyFromObject(val)
	case rbacv1.ClusterRole:
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

// ApplyFromFile applies clusterrole from yaml or json file.
func (h *Handler) ApplyFromFile(filename string) (cr *rbacv1.ClusterRole, err error) {
	cr, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if clusterrole already exist, update it.
		cr, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply clusterrole from bytes data.
func (h *Handler) ApplyFromBytes(data []byte) (cr *rbacv1.ClusterRole, err error) {
	cr, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		cr, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies clusterrole from metav1.Object or runtime.Object.
func (h *Handler) ApplyFromObject(obj interface{}) (*rbacv1.ClusterRole, error) {
	cr, ok := obj.(*rbacv1.ClusterRole)
	if !ok {
		return nil, fmt.Errorf("object type is not *rbacv1.ClusterRole")
	}
	return h.applyCR(cr)
}

// ApplyFromUnstructured applies clusterrole from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*rbacv1.ClusterRole, error) {
	cr := &rbacv1.ClusterRole{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cr)
	if err != nil {
		return nil, err
	}
	return h.applyCR(cr)
}

// ApplyFromMap applies clusterrole from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*rbacv1.ClusterRole, error) {
	cr := &rbacv1.ClusterRole{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cr)
	if err != nil {
		return nil, err
	}
	return h.applyCR(cr)
}

// applyCR
func (h *Handler) applyCR(cr *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
	_, err := h.createCR(cr)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateCR(cr)
	}
	return cr, err
}
