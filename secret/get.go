package secret

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get secret from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.Secret, error) {
	secretJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{}
	err = json.Unmarshal(secretJson, secret)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(secret.Namespace) != 0 {
		namespace = secret.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).GetByName(secret.Name)
}

// GetFromFile get secret from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.Secret, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get secret by name.
func (h *Handler) GetByName(name string) (*corev1.Secret, error) {
	return h.clientset.CoreV1().Secrets(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get secret by name, alias to "GetByName".
func (h *Handler) Get(name string) (*corev1.Secret, error) {
	return h.GetByName(name)
}
