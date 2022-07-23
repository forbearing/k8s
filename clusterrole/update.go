package clusterrole

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates clusterrole from type string, []byte, *rbacv1.ClusterRole,
// rbacv1.ClusterRole, runtime.Object or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*rbacv1.ClusterRole, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *rbacv1.ClusterRole:
		return h.UpdateFromObject(val)
	case rbacv1.ClusterRole:
		return h.UpdateFromObject(&val)
	case runtime.Object:
		return h.UpdateFromObject(val)
	case map[string]interface{}:
		return h.UpdateFromUnstructured(val)
	default:
		return nil, ERR_TYPE_UPDATE
	}
}

// UpdateFromFile updates clusterrole from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*rbacv1.ClusterRole, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates clusterrole from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*rbacv1.ClusterRole, error) {
	crJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cr := &rbacv1.ClusterRole{}
	err = json.Unmarshal(crJson, cr)
	if err != nil {
		return nil, err
	}
	return h.updateClusterRole(cr)
}

// UpdateFromObject updates clusterrole from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*rbacv1.ClusterRole, error) {
	cr, ok := obj.(*rbacv1.ClusterRole)
	if !ok {
		return nil, fmt.Errorf("object is not *rbacv1.ClusterRole")
	}
	return h.updateClusterRole(cr)
}

// UpdateFromUnstructured updates clusterrole from map[string]interface{}.
func (h *Handler) UpdateFromUnstructured(u map[string]interface{}) (*rbacv1.ClusterRole, error) {
	cr := &rbacv1.ClusterRole{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cr)
	if err != nil {
		return nil, err
	}
	return h.updateClusterRole(cr)
}

// updateClusterRole
func (h *Handler) updateClusterRole(cr *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
	cr.ResourceVersion = ""
	cr.UID = ""
	return h.clientset.RbacV1().ClusterRoles().Update(h.ctx, cr, h.Options.UpdateOptions)
}
