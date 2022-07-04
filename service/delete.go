package service

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete service from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	serviceJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	service := &corev1.Service{}
	err = json.Unmarshal(serviceJson, service)
	if err != nil {
		return err
	}

	var namespace string
	if len(service.Namespace) != 0 {
		namespace = service.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(service.Name)
}

// DeleteFromFile delete service from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete service by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().Services(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete service by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
