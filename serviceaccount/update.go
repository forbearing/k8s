package serviceaccount

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update serviceaccount from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*corev1.ServiceAccount, error) {
	sa := &corev1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, sa)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(sa.Namespace) != 0 {
		namespace = sa.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.CoreV1().ServiceAccounts(namespace).Update(h.ctx, sa, h.Options.UpdateOptions)
}

// UpdateFromBytes update serviceaccount from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.ServiceAccount, error) {
	saJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sa := &corev1.ServiceAccount{}
	err = json.Unmarshal(saJson, sa)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(sa.Namespace) != 0 {
		namespace = sa.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.CoreV1().ServiceAccounts(namespace).Update(h.ctx, sa, h.Options.UpdateOptions)
}

// UpdateFromFile update serviceaccount from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.ServiceAccount, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update serviceaccount from yaml file, alias to "UpdateFromFile"
func (h *Handler) Update(filename string) (*corev1.ServiceAccount, error) {
	return h.UpdateFromFile(filename)
}
