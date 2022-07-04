package service

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create service from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*corev1.Service, error) {
	service := &corev1.Service{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, service)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(service.Namespace) != 0 {
		namespace = service.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.CoreV1().Services(namespace).Create(h.ctx, service, h.Options.CreateOptions)
}

// CreateFromBytes create service from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.Service, error) {
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

	return h.clientset.CoreV1().Services(namespace).Create(h.ctx, service, h.Options.CreateOptions)
}

// CreateFromFile create service from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.Service, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create service from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*corev1.Service, error) {
	return h.CreateFromFile(filename)
}
