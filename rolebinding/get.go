package rolebinding

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get rolebinding from bytes.
func (h *Handler) GetFromBytes(data []byte) (*rbacv1.RoleBinding, error) {
	rolebindingJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rolebinding := &rbacv1.RoleBinding{}
	err = json.Unmarshal(rolebindingJson, rolebinding)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(rolebinding.Namespace) != 0 {
		namespace = rolebinding.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).GetByName(rolebinding.Name)
}

// GetFromFile get rolebinding from yaml file.
func (h *Handler) GetFromFile(filename string) (*rbacv1.RoleBinding, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get rolebinding by name.
func (h *Handler) GetByName(name string) (*rbacv1.RoleBinding, error) {
	return h.clientset.RbacV1().RoleBindings(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}
func (h *Handler) Get(name string) (*rbacv1.RoleBinding, error) {
	return h.GetByName(name)
}
