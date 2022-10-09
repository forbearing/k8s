package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates service from type string, []byte, *corev1.Service,
// corev1.Service, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*corev1.Service, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *corev1.Service:
		return h.UpdateFromObject(val)
	case corev1.Service:
		return h.UpdateFromObject(&val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.UpdateFromObject(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates service from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.Service, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates service from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.Service, error) {
	svcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	svc := &corev1.Service{}
	if err = json.Unmarshal(svcJson, svc); err != nil {
		return nil, err
	}
	return h.updateService(svc)
}

// UpdateFromObject updates service from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*corev1.Service, error) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Service")
	}
	return h.updateService(svc)
}

// UpdateFromUnstructured updates service from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*corev1.Service, error) {
	svc := &corev1.Service{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), svc)
	if err != nil {
		return nil, err
	}
	return h.updateService(svc)
}

// UpdateFromMap updates service from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*corev1.Service, error) {
	svc := &corev1.Service{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, svc)
	if err != nil {
		return nil, err
	}
	return h.updateService(svc)
}

// updateService
func (h *Handler) updateService(svc *corev1.Service) (*corev1.Service, error) {
	namespace := svc.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	svc.ResourceVersion = ""
	svc.UID = ""
	return h.clientset.CoreV1().Services(namespace).Update(h.ctx, svc, h.Options.UpdateOptions)
}
