package serviceaccount

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create serviceaccount from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*corev1.ServiceAccount, error) {
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

	return h.clientset.CoreV1().ServiceAccounts(namespace).Create(h.ctx, sa, h.Options.CreateOptions)
}

// CreateFromBytes create serviceaccount from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.ServiceAccount, error) {
	saJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sa := &corev1.ServiceAccount{}
	if err = json.Unmarshal(saJson, sa); err != nil {
		return nil, err
	}

	var namespace string
	if len(sa.Namespace) != 0 {
		namespace = sa.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.CoreV1().ServiceAccounts(namespace).Create(h.ctx, sa, h.Options.CreateOptions)
}

// CreateFromFile create serviceaccount from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.ServiceAccount, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create serviceaccount from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*corev1.ServiceAccount, error) {
	return h.CreateFromFile(filename)
}
