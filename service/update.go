package service

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update service from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*corev1.Service, error) {
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

	return h.clientset.CoreV1().Services(namespace).Update(h.ctx, service, h.Options.UpdateOptions)
}

// UpdateFromBytes update service from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.Service, error) {
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

	return h.clientset.CoreV1().Services(namespace).Update(h.ctx, service, h.Options.UpdateOptions)
}

// UpdateFromFile update service from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.Service, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update service from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*corev1.Service, error) {
	return h.UpdateFromFile(filename)
}
