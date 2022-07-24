package role

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates role from type string, []byte, *rbacv1.Role,
// rbacv1.Role, runtime.Object or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*rbacv1.Role, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *rbacv1.Role:
		return h.CreateFromObject(val)
	case rbacv1.Role:
		return h.CreateFromObject(&val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case map[string]interface{}:
		return h.CreateFromUnstructured(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates role from yaml file.
func (h *Handler) CreateFromFile(filename string) (*rbacv1.Role, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates role from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*rbacv1.Role, error) {
	roleJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	role := &rbacv1.Role{}
	err = json.Unmarshal(roleJson, role)
	if err != nil {
		return nil, err
	}
	return h.createRole(role)
}

// CreateFromObject creates role from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*rbacv1.Role, error) {
	role, ok := obj.(*rbacv1.Role)
	if !ok {
		return nil, fmt.Errorf("object is not *rbacv1.Role")
	}
	return h.createRole(role)
}

// CreateFromUnstructured creates role from map[string]interface{}.
func (h *Handler) CreateFromUnstructured(u map[string]interface{}) (*rbacv1.Role, error) {
	role := &rbacv1.Role{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, role)
	if err != nil {
		return nil, err
	}
	return h.createRole(role)
}

// createRole
func (h *Handler) createRole(role *rbacv1.Role) (*rbacv1.Role, error) {
	var namespace string
	if len(role.Namespace) != 0 {
		namespace = role.Namespace
	} else {
		namespace = h.namespace
	}
	role.ResourceVersion = ""
	role.UID = ""
	return h.clientset.RbacV1().Roles(namespace).Create(h.ctx, role, h.Options.CreateOptions)
}
