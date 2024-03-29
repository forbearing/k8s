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

// Delete deletes role from type string, []byte, *rbacv1.Role,
// rbacv1.Role, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a role from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *rbacv1.Role:
		return h.DeleteFromObject(val)
	case rbacv1.Role:
		return h.DeleteFromObject(&val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	case metav1.Object, runtime.Object:
		return h.DeleteFromObject(val)
	default:
		return ErrInvalidDeleteType
	}
}

// DeleteByName deletes role by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.RbacV1().Roles(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes role from yaml or json file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes role from bytes data.
func (h *Handler) DeleteFromBytes(data []byte) error {
	roleJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	role := &rbacv1.Role{}
	if err = json.Unmarshal(roleJson, role); err != nil {
		return err
	}
	return h.deleteRole(role)
}

// DeleteFromObject deletes role from metav1.Object or runtime.Object.
func (h *Handler) DeleteFromObject(obj interface{}) error {
	role, ok := obj.(*rbacv1.Role)
	if !ok {
		return fmt.Errorf("object type is not *rbacv1.Role")
	}
	return h.deleteRole(role)
}

// DeleteFromUnstructured deletes role from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	role := &rbacv1.Role{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), role)
	if err != nil {
		return err
	}
	return h.deleteRole(role)
}

// DeleteFromMap deletes role from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	role := &rbacv1.Role{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, role)
	if err != nil {
		return err
	}
	return h.deleteRole(role)
}

// deleteRole
func (h *Handler) deleteRole(role *rbacv1.Role) error {
	namespace := role.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.RbacV1().Roles(namespace).Delete(h.ctx, role.Name, h.Options.DeleteOptions)
}
