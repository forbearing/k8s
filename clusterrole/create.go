package clusterrole

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create clusterrole from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*rbacv1.ClusterRole, error) {
	cr := &rbacv1.ClusterRole{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, cr)
	if err != nil {
		return nil, err
	}

	return h.clientset.RbacV1().ClusterRoles().Create(h.ctx, cr, h.Options.CreateOptions)
}

// CreateFromBytes create clusterrole from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*rbacv1.ClusterRole, error) {
	crJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cr := &rbacv1.ClusterRole{}
	if err = json.Unmarshal(crJson, cr); err != nil {
		return nil, err
	}

	return h.clientset.RbacV1().ClusterRoles().Create(h.ctx, cr, h.Options.CreateOptions)
}

// CreateFromFile create clusterrole from yaml file.
func (h *Handler) CreateFromFile(filename string) (*rbacv1.ClusterRole, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create clusterrole from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*rbacv1.ClusterRole, error) {
	return h.CreateFromFile(filename)
}
