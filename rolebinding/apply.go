package rolebinding

import (
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies rolebinding from type string, []byte, *rbacv1.RoleBinding,
// rbacv1.RoleBinding, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*rbacv1.RoleBinding, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *rbacv1.RoleBinding:
		return h.ApplyFromObject(val)
	case rbacv1.RoleBinding:
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

// ApplyFromFile applies rolebinding from yaml file.
func (h *Handler) ApplyFromFile(filename string) (rb *rbacv1.RoleBinding, err error) {
	rb, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if rolebinding already exist, update it.
		rb, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply rolebinding from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (rb *rbacv1.RoleBinding, err error) {
	rb, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		rb, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies rolebinding from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*rbacv1.RoleBinding, error) {
	rb, ok := obj.(*rbacv1.RoleBinding)
	if !ok {
		return nil, fmt.Errorf("object type is not *rbacv1.RoleBinding")
	}
	return h.applyRolebinding(rb)
}

// ApplyFromUnstructured applies rolebinding from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*rbacv1.RoleBinding, error) {
	rb := &rbacv1.RoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rb)
	if err != nil {
		return nil, err
	}
	return h.applyRolebinding(rb)
}

// ApplyFromMap applies rolebinding from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*rbacv1.RoleBinding, error) {
	rb := &rbacv1.RoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rb)
	if err != nil {
		return nil, err
	}
	return h.applyRolebinding(rb)
}

// applyRolebinding
func (h *Handler) applyRolebinding(rb *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error) {
	_, err := h.createRolebinding(rb)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateRolebinding(rb)
	}
	return rb, err
}
