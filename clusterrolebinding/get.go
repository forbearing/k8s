package clusterrolebinding

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets clusterrolebinding from type string, []byte, *rbacv1.ClusterRoleBinding,
// rbacv1.ClusterRoleBinding, runtime.Object or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a clusterrolebinding from file path.
func (h *Handler) Get(obj interface{}) (*rbacv1.ClusterRoleBinding, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *rbacv1.ClusterRoleBinding:
		return h.GetFromObject(val)
	case rbacv1.ClusterRoleBinding:
		return h.GetFromObject(&val)
	case map[string]interface{}:
		return h.GetFromUnstructured(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets clusterrolebinding by name.
func (h *Handler) GetByName(name string) (*rbacv1.ClusterRoleBinding, error) {
	return h.clientset.RbacV1().ClusterRoleBindings().Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets clusterrolebinding from yaml file.
func (h *Handler) GetFromFile(filename string) (*rbacv1.ClusterRoleBinding, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets clusterrolebinding from bytes.
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
	return h.getCRB(crb)
}

// GetFromObject gets clusterrolebinding from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*rbacv1.ClusterRoleBinding, error) {
	crb, ok := obj.(*rbacv1.ClusterRoleBinding)
	if !ok {
		return nil, fmt.Errorf("object is not *rbacv1.ClusterRoleBinding")
	}
	return h.getCRB(crb)
}

// GetFromUnstructured gets clusterrolebinding from map[string]interface{}.
func (h *Handler) GetFromUnstructured(u map[string]interface{}) (*rbacv1.ClusterRoleBinding, error) {
	crb := &rbacv1.ClusterRoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, crb)
	if err != nil {
		return nil, err
	}
	return h.getCRB(crb)
}

// getCRB
// It's necessary to get a new clusterrolebinding resource from a old clusterrolebinding resource,
// because old clusterrolebinding usually don't have clusterrolebinding.Status field.
func (h *Handler) getCRB(crb *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
	return h.clientset.RbacV1().ClusterRoleBindings().Get(h.ctx, crb.Name, h.Options.GetOptions)
}
