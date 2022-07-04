package secret

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update secret from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*corev1.Secret, error) {
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

	return h.clientset.CoreV1().Secrets(namespace).Update(h.ctx, secret, h.Options.UpdateOptions)
}

// UpdateFromBytes update secret from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.Secret, error) {
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

	return h.clientset.CoreV1().Secrets(namespace).Update(h.ctx, secret, h.Options.UpdateOptions)
}

// UpdateFromFile update secret from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.Secret, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update secret from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*corev1.Secret, error) {
	return h.UpdateFromFile(filename)
}
