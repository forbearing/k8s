package clusterrole

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates clusterrole from type string, []byte, *rbacv1.ClusterRole,
// rbacv1.ClusterRole, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*rbacv1.ClusterRole, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *rbacv1.ClusterRole:
		return h.CreateFromObject(val)
	case rbacv1.ClusterRole:
		return h.CreateFromObject(&val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates clusterrole from yaml file.
func (h *Handler) CreateFromFile(filename string) (*rbacv1.ClusterRole, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates clusterrole from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*rbacv1.ClusterRole, error) {
	crJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cr := &rbacv1.ClusterRole{}
	if err = json.Unmarshal(crJson, cr); err != nil {
		return nil, err
	}
	return h.createCR(cr)
}

// CreateFromObject creates clusterrole from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*rbacv1.ClusterRole, error) {
	cr, ok := obj.(*rbacv1.ClusterRole)
	if !ok {
		return nil, fmt.Errorf("object type is not *rbacv1.ClusterRole")
	}
	return h.createCR(cr)
}

// CreateFromUnstructured creates clusterrole from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*rbacv1.ClusterRole, error) {
	cr := &rbacv1.ClusterRole{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cr)
	if err != nil {
		return nil, err
	}
	return h.createCR(cr)
}

// CreateFromMap creates clusterrole from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*rbacv1.ClusterRole, error) {
	cr := &rbacv1.ClusterRole{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cr)
	if err != nil {
		return nil, err
	}
	return h.createCR(cr)
}

// createCR
func (h *Handler) createCR(cr *rbacv1.ClusterRole) (*rbacv1.ClusterRole, error) {
	cr.ResourceVersion = ""
	cr.UID = ""
	return h.clientset.RbacV1().ClusterRoles().Create(h.ctx, cr, h.Options.CreateOptions)
}
