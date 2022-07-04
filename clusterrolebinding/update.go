package clusterrolebinding

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update clusterrolebinding from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*rbacv1.ClusterRoleBinding, error) {
	crb := &rbacv1.ClusterRoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, crb)
	if err != nil {
		return nil, err
	}

	return h.clientset.RbacV1().ClusterRoleBindings().Update(h.ctx, crb, h.Options.UpdateOptions)
}

// UpdateFromBytes update clusterrolebinding from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*rbacv1.ClusterRoleBinding, error) {
	crbJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	crb := &rbacv1.ClusterRoleBinding{}
	err = json.Unmarshal(crbJson, crb)
	if err != nil {
		return nil, err
	}

	return h.clientset.RbacV1().ClusterRoleBindings().Update(h.ctx, crb, h.Options.UpdateOptions)
}

// UpdateFromFile update clusterrolebinding from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*rbacv1.ClusterRoleBinding, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update clusterrolebinding from file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*rbacv1.ClusterRoleBinding, error) {
	return h.UpdateFromFile(filename)
}
