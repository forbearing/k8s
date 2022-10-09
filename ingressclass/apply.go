package ingressclass

import (
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies ingressclass from type string, []byte, *networkingv1.IngressClass,
// networkingv1.IngressClass, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*networkingv1.IngressClass, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *networkingv1.IngressClass:
		return h.ApplyFromObject(val)
	case networkingv1.IngressClass:
		return h.ApplyFromObject(&val)
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	case metav1.Object, runtime.Object:
		return h.ApplyFromObject(val)
	default:
		return nil, ErrInvalidApplyType
	}
}

// ApplyFromFile applies ingressclass from yaml or json file.
func (h *Handler) ApplyFromFile(filename string) (ingc *networkingv1.IngressClass, err error) {
	ingc, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if ingressclass already exist, update it.
		ingc, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply ingressclass from bytes data.
func (h *Handler) ApplyFromBytes(data []byte) (ingc *networkingv1.IngressClass, err error) {
	ingc, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		ingc, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies ingressclass from metav1.Object or runtime.Object.
func (h *Handler) ApplyFromObject(obj interface{}) (*networkingv1.IngressClass, error) {
	ingc, ok := obj.(*networkingv1.IngressClass)
	if !ok {
		return nil, fmt.Errorf("object type is not *networkingv1.IngressClass")
	}
	return h.applyIngressclass(ingc)
}

// ApplyFromUnstructured applies ingressclass from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*networkingv1.IngressClass, error) {
	ingc := &networkingv1.IngressClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ingc)
	if err != nil {
		return nil, err
	}
	return h.applyIngressclass(ingc)
}

// ApplyFromMap applies ingressclass from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*networkingv1.IngressClass, error) {
	ingc := &networkingv1.IngressClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ingc)
	if err != nil {
		return nil, err
	}
	return h.applyIngressclass(ingc)
}

// applyIngressclass
func (h *Handler) applyIngressclass(ingc *networkingv1.IngressClass) (*networkingv1.IngressClass, error) {
	_, err := h.createIngressclass(ingc)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateIngressclass(ingc)
	}
	return ingc, err
}
