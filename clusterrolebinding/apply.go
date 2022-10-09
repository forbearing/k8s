package clusterrolebinding

import (
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies clusterrolebinding from type string, []byte,
// *rbacv1.ClusterRoleBinding, rbacv1.ClusterRoleBinding, metav1.Object, runtime.Object,
// *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*rbacv1.ClusterRoleBinding, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *rbacv1.ClusterRoleBinding:
		return h.ApplyFromObject(val)
	case rbacv1.ClusterRoleBinding:
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

// ApplyFromFile applies clusterrolebinding from yaml or json file.
func (h *Handler) ApplyFromFile(filename string) (crb *rbacv1.ClusterRoleBinding, err error) {
	crb, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if clusterrolebinding already exist, update it.
		crb, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply clusterrolebinding from bytes data.
func (h *Handler) ApplyFromBytes(data []byte) (crb *rbacv1.ClusterRoleBinding, err error) {
	crb, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		crb, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies clusterrolebinding from metav1.Object or runtime.Object.
func (h *Handler) ApplyFromObject(obj interface{}) (*rbacv1.ClusterRoleBinding, error) {
	crb, ok := obj.(*rbacv1.ClusterRoleBinding)
	if !ok {
		return nil, fmt.Errorf("object type is not *rbacv1.ClusterRoleBinding")
	}
	return h.applyCRB(crb)
}

// ApplyFromUnstructured applies clusterrolebinding from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*rbacv1.ClusterRoleBinding, error) {
	crb := &rbacv1.ClusterRoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), crb)
	if err != nil {
		return nil, err
	}
	return h.applyCRB(crb)
}

// ApplyFromMap applies clusterrolebinding from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*rbacv1.ClusterRoleBinding, error) {
	crb := &rbacv1.ClusterRoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, crb)
	if err != nil {
		return nil, err
	}
	return h.applyCRB(crb)
}

// applyCRB
func (h *Handler) applyCRB(crb *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
	_, err := h.createCRB(crb)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateCRB(crb)
	}
	return crb, err
}
