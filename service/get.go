package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets service from type string, []byte, *corev1.Service,
// corev1.Service, runtime.Object or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a service from file path.
func (h *Handler) Get(obj interface{}) (*corev1.Service, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *corev1.Service:
		return h.GetFromObject(val)
	case corev1.Service:
		return h.GetFromObject(&val)
	case runtime.Object:
		return h.GetFromObject(val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets service by name.
func (h *Handler) GetByName(name string) (*corev1.Service, error) {
	return h.clientset.CoreV1().Services(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets service from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.Service, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets service from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.Service, error) {
	svcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	svc := &corev1.Service{}
	if err = json.Unmarshal(svcJson, svc); err != nil {
		return nil, err
	}
	return h.getService(svc)
}

// GetFromObject gets service from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*corev1.Service, error) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Service")
	}
	return h.getService(svc)
}

// GetFromUnstructured gets service from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*corev1.Service, error) {
	svc := &corev1.Service{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), svc)
	if err != nil {
		return nil, err
	}
	return h.getService(svc)
}

// GetFromMap gets service from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*corev1.Service, error) {
	svc := &corev1.Service{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, svc)
	if err != nil {
		return nil, err
	}
	return h.getService(svc)
}

// getService
// It's necessary to get a new service resource from a old service resource,
// because old service usually don't have service.Status field.
func (h *Handler) getService(svc *corev1.Service) (*corev1.Service, error) {
	var namespace string
	if len(svc.Namespace) != 0 {
		namespace = svc.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().Services(namespace).Get(h.ctx, svc.Name, h.Options.GetOptions)
}
