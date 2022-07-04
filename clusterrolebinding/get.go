package clusterrolebinding

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get clusterrolebinding from bytes.
func (h *Handler) GetFromBytes(data []byte) (*rbacv1.ClusterRoleBinding, error) {
	crbJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	crb := &rbacv1.ClusterRoleBinding{}
	err = json.Unmarshal(crbJson, crb)
	if err != nil {
		return nil, err
	}

	return h.GetByName(crb.Name)
}

// GetFromFile get clusterrolebinding from yaml file.
func (h *Handler) GetFromFile(filename string) (*rbacv1.ClusterRoleBinding, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get clusterrolebinding by name.
func (h *Handler) GetByName(name string) (*rbacv1.ClusterRoleBinding, error) {
	return h.clientset.RbacV1().ClusterRoleBindings().Get(h.ctx, name, h.Options.GetOptions)
}

// Get get clusterrolebinding by name, alias to "GetByName".
func (h *Handler) Get(name string) (*rbacv1.ClusterRoleBinding, error) {
	return h.clientset.RbacV1().ClusterRoleBindings().Get(h.ctx, name, h.Options.GetOptions)
}
