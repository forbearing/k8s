package service

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get service from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.Service, error) {
	serviceJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	service := &corev1.Service{}
	err = json.Unmarshal(serviceJson, service)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(service.Namespace) != 0 {
		namespace = service.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).GetByName(service.Name)
}

// GetFromFile get service from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.Service, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get service by name.
func (h *Handler) GetByName(name string) (*corev1.Service, error) {
	return h.clientset.CoreV1().Services(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get service by name, alias to "GetByName".
func (h *Handler) Get(name string) (*corev1.Service, error) {
	return h.GetByName(name)
}
