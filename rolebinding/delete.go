package rolebinding

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete rolebinding from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	rolebindingJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	rolebinding := &rbacv1.RoleBinding{}
	if err = json.Unmarshal(rolebindingJson, rolebinding); err != nil {
		return err
	}

	var namespace string
	if len(rolebinding.Namespace) != 0 {
		namespace = rolebinding.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(rolebinding.Name)
}

// DeleteFromFile delete rolebinding from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete rolebinding by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.RbacV1().RoleBindings(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete rolebinding by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) (err error) {
	return h.DeleteByName(name)
}
