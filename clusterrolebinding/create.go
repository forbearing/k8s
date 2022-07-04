package clusterrolebinding

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create ClusterRoleBinding from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*rbacv1.ClusterRoleBinding, error) {
	crb := &rbacv1.ClusterRoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, crb)
	if err != nil {
		return nil, err
	}

	return h.clientset.RbacV1().ClusterRoleBindings().Create(h.ctx, crb, h.Options.CreateOptions)
}

// CreateFromBytes create clusterrolebinding from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*rbacv1.ClusterRoleBinding, error) {
	crbJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	crb := &rbacv1.ClusterRoleBinding{}
	err = json.Unmarshal(crbJson, crb)
	if err != nil {
		return nil, err
	}

	return h.clientset.RbacV1().ClusterRoleBindings().Create(h.ctx, crb, h.Options.CreateOptions)
}

// CreateFromFile create clusterrolebinding from yaml file.
func (h *Handler) CreateFromFile(filename string) (*rbacv1.ClusterRoleBinding, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create clusterrolebinding from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*rbacv1.ClusterRoleBinding, error) {
	return h.CreateFromFile(filename)
}
