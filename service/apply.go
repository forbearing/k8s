package service

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies service from type string, []byte, *corev1.Service,
// corev1.Service, runtime.Object or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*corev1.Service, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *corev1.Service:
		return h.ApplyFromObject(val)
	case corev1.Service:
		return h.ApplyFromObject(&val)
	case runtime.Object:
		return h.ApplyFromObject(val)
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	default:
		return nil, ERR_TYPE_APPLY
	}
}

// ApplyFromFile applies service from yaml file.
func (h *Handler) ApplyFromFile(filename string) (svc *corev1.Service, err error) {
	svc, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if service already exist, update it.
		svc, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply service from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (svc *corev1.Service, err error) {
	svc, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		svc, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies service from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*corev1.Service, error) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Service")
	}
	return h.applyService(svc)
}

// ApplyFromUnstructured applies service from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*corev1.Service, error) {
	svc := &corev1.Service{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), svc)
	if err != nil {
		return nil, err
	}
	return h.applyService(svc)
}

// ApplyFromMap applies service from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*corev1.Service, error) {
	svc := &corev1.Service{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, svc)
	if err != nil {
		return nil, err
	}
	return h.applyService(svc)
}

// applyService
func (h *Handler) applyService(svc *corev1.Service) (*corev1.Service, error) {
	_, err := h.createService(svc)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateService(svc)
	}
	return svc, err
}
