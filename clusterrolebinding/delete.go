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

// Delete deletes clusterrolebinding from type string, []byte,
// *rbacv1.ClusterRoleBinding, rbacv1.ClusterRoleBinding, runtime.Object,
// *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a clusterrolebinding from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *rbacv1.ClusterRoleBinding:
		return h.DeleteFromObject(val)
	case rbacv1.ClusterRoleBinding:
		return h.DeleteFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.DeleteFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.DeleteFromObject(val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	default:
		return ERR_TYPE_DELETE
	}
}

// DeleteByName deletes clusterrolebinding by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.RbacV1().ClusterRoleBindings().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes clusterrolebinding from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes clusterrolebinding from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	crbJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	crb := &rbacv1.ClusterRoleBinding{}
	if err = json.Unmarshal(crbJson, crb); err != nil {
		return err
	}
	return h.deleteCRB(crb)
}

// DeleteFromObject deletes clusterrolebinding from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	crb, ok := obj.(*rbacv1.ClusterRoleBinding)
	if !ok {
		return fmt.Errorf("object type is not *rbacv1.ClusterRoleBinding")
	}
	return h.deleteCRB(crb)
}

// DeleteFromUnstructured deletes clusterrolebinding from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	crb := &rbacv1.ClusterRoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), crb)
	if err != nil {
		return err
	}
	return h.deleteCRB(crb)
}

// DeleteFromMap deletes clusterrolebinding from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	crb := &rbacv1.ClusterRoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, crb)
	if err != nil {
		return err
	}
	return h.deleteCRB(crb)
}

// deleteCRB
func (h *Handler) deleteCRB(crb *rbacv1.ClusterRoleBinding) error {
	return h.clientset.RbacV1().ClusterRoleBindings().Delete(h.ctx, crb.Name, h.Options.DeleteOptions)
}
