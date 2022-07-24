package rolebinding

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets rolebinding from type string, []byte, *rbacv1.RoleBinding,
// rbacv1.RoleBinding, runtime.Object or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a rolebinding from file path.
func (h *Handler) Get(obj interface{}) (*rbacv1.RoleBinding, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *rbacv1.RoleBinding:
		return h.GetFromObject(val)
	case rbacv1.RoleBinding:
		return h.GetFromObject(&val)
	case map[string]interface{}:
		return h.GetFromUnstructured(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets rolebinding by name.
func (h *Handler) GetByName(name string) (*rbacv1.RoleBinding, error) {
	return h.clientset.RbacV1().RoleBindings(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets rolebinding from yaml file.
func (h *Handler) GetFromFile(filename string) (*rbacv1.RoleBinding, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets rolebinding from bytes.
func (h *Handler) GetFromBytes(data []byte) (*rbacv1.RoleBinding, error) {
	rbJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rb := &rbacv1.RoleBinding{}
	err = json.Unmarshal(rbJson, rb)
	if err != nil {
		return nil, err
	}
	return h.getRolebinding(rb)
}

// GetFromObject gets rolebinding from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*rbacv1.RoleBinding, error) {
	rb, ok := obj.(*rbacv1.RoleBinding)
	if !ok {
		return nil, fmt.Errorf("object is not *rbacv1.RoleBinding")
	}
	return h.getRolebinding(rb)
}

// GetFromUnstructured gets rolebinding from map[string]interface{}.
func (h *Handler) GetFromUnstructured(u map[string]interface{}) (*rbacv1.RoleBinding, error) {
	rb := &rbacv1.RoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rb)
	if err != nil {
		return nil, err
	}
	return h.getRolebinding(rb)
}

// getRolebinding
// It's necessary to get a new rolebinding resource from a old rolebinding resource,
// because old rolebinding usually don't have rolebinding.Status field.
func (h *Handler) getRolebinding(rb *rbacv1.RoleBinding) (*rbacv1.RoleBinding, error) {
	var namespace string
	if len(rb.Namespace) != 0 {
		namespace = rb.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.RbacV1().RoleBindings(namespace).Get(h.ctx, rb.Name, h.Options.GetOptions)
}
