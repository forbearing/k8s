package rolebinding

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

// Create creates rolebinding from type string, []byte, *rbacv1.RoleBinding,
// rbacv1.RoleBinding, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*rbacv1.RoleBinding, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *rbacv1.RoleBinding:
		return h.CreateFromObject(val)
	case rbacv1.RoleBinding:
		return h.CreateFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.CreateFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.CreateFromObject(val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	default:
		return nil, ErrInvalidCreateType
	}
}

// CreateFromFile creates rolebinding from yaml file.
func (h *Handler) CreateFromFile(filename string) (*rbacv1.RoleBinding, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates rolebinding from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*rbacv1.RoleBinding, error) {
	rbJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rb := &rbacv1.RoleBinding{}
	if err = json.Unmarshal(rbJson, rb); err != nil {
		return nil, err
	}
	return h.createRolebinding(rb)
}

// CreateFromObject creates rolebinding from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*rbacv1.RoleBinding, error) {
	rb, ok := obj.(*rbacv1.RoleBinding)
	if !ok {
		return nil, fmt.Errorf("object type is not *rbacv1.RoleBinding")
	}
	return h.createRolebinding(rb)
}

// CreateFromUnstructured creates rolebinding from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*rbacv1.RoleBinding, error) {
	rb := &rbacv1.RoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rb)
	if err != nil {
		return nil, err
	}
	return h.createRolebinding(rb)
}

// CreateFromMap creates rolebinding from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*rbacv1.RoleBinding, error) {
	rb := &rbacv1.RoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rb)
	if err != nil {
		return nil, err
	}
	return h.createRolebinding(rb)
}

// createRolebinding
func (h *Handler) createRolebinding(rb *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error) {
	var namespace string
	if len(rb.Namespace) != 0 {
		namespace = rb.Namespace
	} else {
		namespace = h.namespace
	}
	rb.ResourceVersion = ""
	rb.UID = ""
	return h.clientset.RbacV1().RoleBindings(namespace).Create(h.ctx, rb, h.Options.CreateOptions)
}
