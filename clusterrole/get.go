package clusterrole

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get clusterrole from bytes.
func (h *Handler) GetFromBytes(data []byte) (*rbacv1.ClusterRole, error) {
	crJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cr := &rbacv1.ClusterRole{}
	err = json.Unmarshal(crJson, cr)
	if err != nil {
		return nil, err
	}

	return h.GetByName(cr.Name)
}

// GetFromFile get clusterrole from yaml file.
func (h *Handler) GetFromFile(filename string) (*rbacv1.ClusterRole, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get clusterrole by name.
func (h *Handler) GetByName(name string) (*rbacv1.ClusterRole, error) {
	return h.clientset.RbacV1().ClusterRoles().Get(h.ctx, name, h.Options.GetOptions)
}

// Get get clusterrole by name, alias to "GetByName".
func (h *Handler) Get(name string) (*rbacv1.ClusterRole, error) {
	return h.GetByName(name)
}
