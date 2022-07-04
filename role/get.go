package role

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get role from bytes.
func (h *Handler) GetFromBytes(data []byte) (*rbacv1.Role, error) {
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

	return h.WithNamespace(namespace).GetByName(role.Name)
}

// GetFromFile get role from yaml file.
func (h *Handler) GetFromFile(filename string) (*rbacv1.Role, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get role by name.
func (h *Handler) GetByName(name string) (*rbacv1.Role, error) {
	return h.clientset.RbacV1().Roles(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get role by name, alias to "GetByName".
func (h *Handler) Get(name string) (*rbacv1.Role, error) {
	return h.GetByName(name)
}
