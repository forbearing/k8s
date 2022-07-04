package role

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update role from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*rbacv1.Role, error) {
	role := &rbacv1.Role{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, role)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(role.Namespace) != 0 {
		namespace = role.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.RbacV1().Roles(namespace).Update(h.ctx, role, h.Options.UpdateOptions)
}

// UpdateFromBytes update role from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*rbacv1.Role, error) {
	roleJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	role := &rbacv1.Role{}
	err = json.Unmarshal(roleJson, role)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(role.Namespace) != 0 {
		namespace = role.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.RbacV1().Roles(namespace).Update(h.ctx, role, h.Options.UpdateOptions)
}

// UpdateFromFile update role from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*rbacv1.Role, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update role from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*rbacv1.Role, error) {
	return h.UpdateFromFile(filename)
}
