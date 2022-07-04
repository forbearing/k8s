package serviceaccount

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get serviceaccount from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.ServiceAccount, error) {
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

	return h.WithNamespace(namespace).GetByName(sa.Name)
}

// GetFromFile get serviceaccount from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.ServiceAccount, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get serviceaccount by name.
func (h *Handler) GetByName(name string) (*corev1.ServiceAccount, error) {
	return h.clientset.CoreV1().ServiceAccounts(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get serviceaccount by name.
func (h *Handler) Get(name string) (*corev1.ServiceAccount, error) {
	return h.GetByName(name)
}
