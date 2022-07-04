package serviceaccount

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete serviceaccount from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	saJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	sa := &corev1.ServiceAccount{}
	err = json.Unmarshal(saJson, sa)
	if err != nil {
		return err
	}

	var namespace string
	if len(sa.Namespace) != 0 {
		namespace = sa.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(sa.Name)
}

// DeleteFromFile delete serviceaccount from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete serviceaccount by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().ServiceAccounts(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete serviceaccount by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) (err error) {
	return h.DeleteByName(name)
}
