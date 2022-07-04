package secret

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create secret from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, secret)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(secret.Namespace) != 0 {
		namespace = secret.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.CoreV1().Secrets(namespace).Create(h.ctx, secret, h.Options.CreateOptions)
}

// CreateFromBytes create secret from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.Secret, error) {
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

	return h.clientset.CoreV1().Secrets(namespace).Create(h.ctx, secret, h.Options.CreateOptions)
}

// CreateFromFile create secret from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.Secret, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create secret from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*corev1.Secret, error) {
	return h.CreateFromFile(filename)
}
