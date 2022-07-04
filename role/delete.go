package role

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete role from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	roleJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	role := &rbacv1.Role{}
	err = json.Unmarshal(roleJson, role)
	if err != nil {
		return err
	}

	var namespace string
	if len(role.Namespace) != 0 {
		namespace = role.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(role.Name)
}

// DeleteFromFile delete role from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete role by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.RbacV1().Roles(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete role by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
