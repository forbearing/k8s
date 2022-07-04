package clusterrolebinding

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete clusterrolebinding from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	crbJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	crb := &rbacv1.ClusterRoleBinding{}
	err = json.Unmarshal(crbJson, crb)
	if err != nil {
		return err
	}

	return h.DeleteByName(crb.Name)
}

// DeleteFromFile delete clusterrolebinding from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete clusterrolebinding by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.RbacV1().ClusterRoleBindings().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete clusterrolebinding by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
