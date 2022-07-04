package rolebinding

import (
	"encoding/json"
	"io/ioutil"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update rolebinding from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*rbacv1.RoleBinding, error) {
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

	return h.clientset.RbacV1().RoleBindings(namespace).Update(h.ctx, rolebinding, h.Options.UpdateOptions)
}

// UpdateFromBytes update rolebinding from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*rbacv1.RoleBinding, error) {
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

	return h.clientset.RbacV1().RoleBindings(namespace).Update(h.ctx, rolebinding, h.Options.UpdateOptions)
}

// UpdateFromFile update rolebinding from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*rbacv1.RoleBinding, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update rolebinding from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*rbacv1.RoleBinding, error) {
	return h.UpdateFromFile(filename)
}
