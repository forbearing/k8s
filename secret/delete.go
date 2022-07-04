package secret

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete secret from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	secretJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	secret := &corev1.Secret{}
	err = json.Unmarshal(secretJson, secret)
	if err != nil {
		return err
	}

	var namespace string
	if len(secret.Namespace) != 0 {
		namespace = secret.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(secret.Name)
}

// DeleteFromFile delete secret from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete secret by name
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().Secrets(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete secret by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
