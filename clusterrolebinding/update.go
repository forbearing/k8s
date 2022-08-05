package clusterrolebinding

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates clusterrolebinding from type string, []byte,
// *rbacv1.ClusterRoleBinding, rbacv1.ClusterRoleBinding, runtime.Object,
// *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*rbacv1.ClusterRoleBinding, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *rbacv1.ClusterRoleBinding:
		return h.UpdateFromObject(val)
	case rbacv1.ClusterRoleBinding:
		return h.UpdateFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.UpdateFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.UpdateFromObject(val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates clusterrolebinding from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*rbacv1.ClusterRoleBinding, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates clusterrolebinding from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*rbacv1.ClusterRoleBinding, error) {
	crbJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	crb := &rbacv1.ClusterRoleBinding{}
	if err = json.Unmarshal(crbJson, crb); err != nil {
		return nil, err
	}
	return h.updateCRB(crb)
}

// UpdateFromObject updates clusterrolebinding from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*rbacv1.ClusterRoleBinding, error) {
	crb, ok := obj.(*rbacv1.ClusterRoleBinding)
	if !ok {
		return nil, fmt.Errorf("object type is not *rbacv1.ClusterRoleBinding")
	}
	return h.updateCRB(crb)
}

// UpdateFromUnstructured updates clusterrolebinding from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*rbacv1.ClusterRoleBinding, error) {
	crb := &rbacv1.ClusterRoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), crb)
	if err != nil {
		return nil, err
	}
	return h.updateCRB(crb)
}

// UpdateFromMap updates clusterrolebinding from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*rbacv1.ClusterRoleBinding, error) {
	crb := &rbacv1.ClusterRoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, crb)
	if err != nil {
		return nil, err
	}
	return h.updateCRB(crb)
}

// updateCRB
func (h *Handler) updateCRB(crb *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
	crb.ResourceVersion = ""
	crb.UID = ""
	return h.clientset.RbacV1().ClusterRoleBindings().Update(h.ctx, crb, h.Options.UpdateOptions)
}
