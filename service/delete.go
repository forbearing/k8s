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

// Delete deletes service from type string, []byte, *corev1.Service,
// corev1.Service, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a service from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *corev1.Service:
		return h.DeleteFromObject(val)
	case corev1.Service:
		return h.DeleteFromObject(&val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	case metav1.Object, runtime.Object:
		return h.DeleteFromObject(val)
	default:
		return ErrInvalidDeleteType
	}
}

// DeleteByName deletes service by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().Services(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes service from yaml or json file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes service from bytes data.
func (h *Handler) DeleteFromBytes(data []byte) error {
	svcJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	svc := &corev1.Service{}
	if err = json.Unmarshal(svcJson, svc); err != nil {
		return err
	}
	return h.deleteService(svc)
}

// DeleteFromObject deletes service from metav1.Object or runtime.Object.
func (h *Handler) DeleteFromObject(obj interface{}) error {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		return fmt.Errorf("object type is not *corev1.Service")
	}
	return h.deleteService(svc)
}

// DeleteFromUnstructured deletes service from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	svc := &corev1.Service{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), svc)
	if err != nil {
		return err
	}
	return h.deleteService(svc)
}

// DeleteFromMap deletes service from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	svc := &corev1.Service{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, svc)
	if err != nil {
		return err
	}
	return h.deleteService(svc)
}

// deleteService
func (h *Handler) deleteService(svc *corev1.Service) error {
	namespace := svc.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().Services(namespace).Delete(h.ctx, svc.Name, h.Options.DeleteOptions)
}
