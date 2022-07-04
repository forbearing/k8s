package rolebinding

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create rolebinding from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*rbacv1.RoleBinding, error) {
	rolebinding := &rbacv1.RoleBinding{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, rolebinding)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(rolebinding.Namespace) != 0 {
		namespace = rolebinding.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.RbacV1().RoleBindings(namespace).Create(h.ctx, rolebinding, h.Options.CreateOptions)
}

// CreateFromBytes create rolebinding from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*rbacv1.RoleBinding, error) {

	rolebindingJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rolebinding := &rbacv1.RoleBinding{}
	err = json.Unmarshal(rolebindingJson, rolebinding)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(rolebinding.Namespace) != 0 {
		namespace = rolebinding.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.RbacV1().RoleBindings(namespace).Create(h.ctx, rolebinding, h.Options.CreateOptions)
}

// CreateFromFile create rolebinding from yaml file.
func (h *Handler) CreateFromFile(filename string) (*rbacv1.RoleBinding, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create rolebinding from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*rbacv1.RoleBinding, error) {
	return h.CreateFromFile(filename)
}
