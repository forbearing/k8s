package role

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets role from type string, []byte, *rbacv1.Role,
// rbacv1.Role, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a role from file path.
func (h *Handler) Get(obj interface{}) (*rbacv1.Role, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *rbacv1.Role:
		return h.GetFromObject(val)
	case rbacv1.Role:
		return h.GetFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.GetFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.GetFromObject(val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	default:
		return nil, ErrInvalidGetType
	}
}

// GetByName gets role by name.
func (h *Handler) GetByName(name string) (*rbacv1.Role, error) {
	return h.clientset.RbacV1().Roles(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets role from yaml file.
func (h *Handler) GetFromFile(filename string) (*rbacv1.Role, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets role from bytes.
func (h *Handler) GetFromBytes(data []byte) (*rbacv1.Role, error) {
	roleJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	role := &rbacv1.Role{}
	if err = json.Unmarshal(roleJson, role); err != nil {
		return nil, err
	}
	return h.getRole(role)
}

// GetFromObject gets role from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*rbacv1.Role, error) {
	role, ok := obj.(*rbacv1.Role)
	if !ok {
		return nil, fmt.Errorf("object type is not *rbacv1.Role")
	}
	return h.getRole(role)
}

// GetFromUnstructured gets role from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*rbacv1.Role, error) {
	role := &rbacv1.Role{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), role)
	if err != nil {
		return nil, err
	}
	return h.getRole(role)
}

// GetFromMap gets role from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*rbacv1.Role, error) {
	role := &rbacv1.Role{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, role)
	if err != nil {
		return nil, err
	}
	return h.getRole(role)
}

// getRole
// It's necessary to get a new role resource from a old role resource,
// because old role usually don't have role.Status field.
func (h *Handler) getRole(role *rbacv1.Role) (*rbacv1.Role, error) {
	var namespace string
	if len(role.Namespace) != 0 {
		namespace = role.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.RbacV1().Roles(namespace).Get(h.ctx, role.Name, h.Options.GetOptions)
}
