package role

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates role from type string, []byte, *rbacv1.Role,
// rbacv1.Role, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*rbacv1.Role, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *rbacv1.Role:
		return h.UpdateFromObject(val)
	case rbacv1.Role:
		return h.UpdateFromObject(&val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.UpdateFromObject(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates role from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*rbacv1.Role, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates role from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*rbacv1.Role, error) {
	roleJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	role := &rbacv1.Role{}
	if err = json.Unmarshal(roleJson, role); err != nil {
		return nil, err
	}
	return h.updateRole(role)
}

// UpdateFromObject updates role from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*rbacv1.Role, error) {
	role, ok := obj.(*rbacv1.Role)
	if !ok {
		return nil, fmt.Errorf("object type is not *rbacv1.Role")
	}
	return h.updateRole(role)
}

// UpdateFromUnstructured updates role from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*rbacv1.Role, error) {
	role := &rbacv1.Role{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), role)
	if err != nil {
		return nil, err
	}
	return h.updateRole(role)
}

// UpdateFromMap updates role from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*rbacv1.Role, error) {
	role := &rbacv1.Role{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, role)
	if err != nil {
		return nil, err
	}
	return h.updateRole(role)
}

// updateRole
func (h *Handler) updateRole(role *rbacv1.Role) (*rbacv1.Role, error) {
	namespace := role.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	role.ResourceVersion = ""
	role.UID = ""
	return h.clientset.RbacV1().Roles(namespace).Update(h.ctx, role, h.Options.UpdateOptions)
}
