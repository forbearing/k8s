package clusterrole

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

// Delete deletes clusterrole from type string, []byte, *rbacv1.ClusterRole,
// rbacv1.ClusterRole, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a clusterrole from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *rbacv1.ClusterRole:
		return h.DeleteFromObject(val)
	case rbacv1.ClusterRole:
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
		return ErrInvalidDeleteType
	}
}

// DeleteByName deletes clusterrole by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.RbacV1().ClusterRoles().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes clusterrole from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes clusterrole from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	crJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	cr := &rbacv1.ClusterRole{}
	if err = json.Unmarshal(crJson, cr); err != nil {
		return err
	}
	return h.deleteCR(cr)
}

// DeleteFromObject deletes clusterrole from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	cr, ok := obj.(*rbacv1.ClusterRole)
	if !ok {
		return fmt.Errorf("object type is not *rbacv1.ClusterRole")
	}
	return h.deleteCR(cr)
}

// DeleteFromUnstructured deletes clusterrole from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	cr := &rbacv1.ClusterRole{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cr)
	if err != nil {
		return err
	}
	return h.deleteCR(cr)
}

// DeleteFromMap deletes clusterrole from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	cr := &rbacv1.ClusterRole{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cr)
	if err != nil {
		return err
	}
	return h.deleteCR(cr)
}

// deleteCR
func (h *Handler) deleteCR(cr *rbacv1.ClusterRole) error {
	return h.clientset.RbacV1().ClusterRoles().Delete(h.ctx, cr.Name, h.Options.DeleteOptions)
}
