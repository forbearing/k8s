package clusterrole

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete clusterrole from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	crJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	cr := &rbacv1.ClusterRole{}
	err = json.Unmarshal(crJson, cr)
	if err != nil {
		return err
	}

	return h.DeleteByName(cr.Name)
}

// DeleteFromFile delete clusterrole from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete clusterrole by name.
func (h *Handler) DeleteByName(name string) (err error) {
	return h.clientset.RbacV1().ClusterRoles().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete clusterrole by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) (err error) {
	return h.DeleteByName(name)
}
