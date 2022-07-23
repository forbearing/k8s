package clusterrole

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets clusterrole from type string, []byte, *rbacv1.ClusterRole,
// rbacv1.ClusterRole, runtime.Object or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a clusterrole from file path.
func (h *Handler) Get(obj interface{}) (*rbacv1.ClusterRole, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *rbacv1.ClusterRole:
		return h.GetFromObject(val)
	case rbacv1.ClusterRole:
		return h.GetFromObject(&val)
	case map[string]interface{}:
		return h.GetFromUnstructured(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets clusterrole by name.
func (h *Handler) GetByName(name string) (*rbacv1.ClusterRole, error) {
	return h.clientset.RbacV1().ClusterRoles().Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets clusterrole from yaml file.
func (h *Handler) GetFromFile(filename string) (*rbacv1.ClusterRole, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets clusterrole from bytes.
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
	return h.getClusterRole(cr)
}

// GetFromObject gets clusterrole from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*rbacv1.ClusterRole, error) {
	cr, ok := obj.(*rbacv1.ClusterRole)
	if !ok {
		return nil, fmt.Errorf("object is not *rbacv1.ClusterRole")
	}
	return h.getClusterRole(cr)
}

// GetFromUnstructured gets clusterrole from map[string]interface{}.
func (h *Handler) GetFromUnstructured(u map[string]interface{}) (*rbacv1.ClusterRole, error) {
	cr := &rbacv1.ClusterRole{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cr)
	if err != nil {
		return nil, err
	}
	return h.getClusterRole(cr)
}

// getClusterRole
// It's necessary to get a new clusterrole resource from a old clusterrole resource,
// because old clusterrole usually don't have clusterrole.Status field.
func (h *Handler) getClusterRole(cr *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
	return h.clientset.RbacV1().ClusterRoles().Get(h.ctx, cr.Name, h.Options.GetOptions)
}
