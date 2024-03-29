package rolebinding

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates rolebinding from type string, []byte, *rbacv1.RoleBinding,
// rbacv1.RoleBinding, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*rbacv1.RoleBinding, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *rbacv1.RoleBinding:
		return h.UpdateFromObject(val)
	case rbacv1.RoleBinding:
		return h.UpdateFromObject(&val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.UpdateFromObject(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates rolebinding from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*rbacv1.RoleBinding, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates rolebinding from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*rbacv1.RoleBinding, error) {
	rbJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rb := &rbacv1.RoleBinding{}
	if err = json.Unmarshal(rbJson, rb); err != nil {
		return nil, err
	}
	return h.updateRolebinding(rb)
}

// UpdateFromObject updates rolebinding from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*rbacv1.RoleBinding, error) {
	rb, ok := obj.(*rbacv1.RoleBinding)
	if !ok {
		return nil, fmt.Errorf("object type is not *rbacv1.RoleBinding")
	}
	return h.updateRolebinding(rb)
}

// UpdateFromUnstructured updates rolebinding from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*rbacv1.RoleBinding, error) {
	rb := &rbacv1.RoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rb)
	if err != nil {
		return nil, err
	}
	return h.updateRolebinding(rb)
}

// UpdateFromMap updates rolebinding from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*rbacv1.RoleBinding, error) {
	rb := &rbacv1.RoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rb)
	if err != nil {
		return nil, err
	}
	return h.updateRolebinding(rb)
}

// updateRolebinding
func (h *Handler) updateRolebinding(rb *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error) {
	namespace := rb.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	rb.ResourceVersion = ""
	rb.UID = ""
	return h.clientset.RbacV1().RoleBindings(namespace).Update(h.ctx, rb, h.Options.UpdateOptions)
}
