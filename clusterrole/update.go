package clusterrole

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update clusterrole from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*rbacv1.ClusterRole, error) {
	cr := &rbacv1.ClusterRole{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, cr)
	if err != nil {
		return nil, err
	}

	return h.clientset.RbacV1().ClusterRoles().Update(h.ctx, cr, h.Options.UpdateOptions)
}

// UpdateFromBytes update clusterrole from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*rbacv1.ClusterRole, error) {
	crJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cr := &rbacv1.ClusterRole{}
	if err = json.Unmarshal(crJson, cr); err != nil {
		return nil, err
	}

	return h.clientset.RbacV1().ClusterRoles().Update(h.ctx, cr, h.Options.UpdateOptions)
}

// UpdateFromFile update clusterrole from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*rbacv1.ClusterRole, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update clusterrole from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*rbacv1.ClusterRole, error) {
	return h.UpdateFromFile(filename)
}
