package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates service from type string, []byte, *corev1.Service,
// corev1.Service, runtime.Object or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*corev1.Service, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *corev1.Service:
		return h.CreateFromObject(val)
	case corev1.Service:
		return h.CreateFromObject(&val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case map[string]interface{}:
		return h.CreateFromUnstructured(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates service from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.Service, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates service from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.Service, error) {
	svcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	svc := &corev1.Service{}
	err = json.Unmarshal(svcJson, svc)
	if err != nil {
		return nil, err
	}
	return h.createService(svc)
}

// CreateFromObject creates service from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*corev1.Service, error) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		return nil, fmt.Errorf("object is not *corev1.Service")
	}
	return h.createService(svc)
}

// CreateFromUnstructured creates service from map[string]interface{}.
func (h *Handler) CreateFromUnstructured(u map[string]interface{}) (*corev1.Service, error) {
	svc := &corev1.Service{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, svc)
	if err != nil {
		return nil, err
	}
	return h.createService(svc)
}

// createService
func (h *Handler) createService(svc *corev1.Service) (*corev1.Service, error) {
	var namespace string
	if len(svc.Namespace) != 0 {
		namespace = svc.Namespace
	} else {
		namespace = h.namespace
	}
	svc.ResourceVersion = ""
	svc.UID = ""
	return h.clientset.CoreV1().Services(namespace).Create(h.ctx, svc, h.Options.CreateOptions)
}
