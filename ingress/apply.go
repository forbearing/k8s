package ingress

import (
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies ingress from type string, []byte, *networkingv1.Ingress,
// networkingv1.Ingress, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*networkingv1.Ingress, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *networkingv1.Ingress:
		return h.ApplyFromObject(val)
	case networkingv1.Ingress:
		return h.ApplyFromObject(&val)
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	case runtime.Object:
		return h.ApplyFromObject(val)
	default:
		return nil, ErrInvalidApplyType
	}
}

// ApplyFromFile applies ingress from yaml file.
func (h *Handler) ApplyFromFile(filename string) (ing *networkingv1.Ingress, err error) {
	ing, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if ingress already exist, update it.
		ing, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply ingress from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (ing *networkingv1.Ingress, err error) {
	ing, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		ing, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies ingress from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*networkingv1.Ingress, error) {
	ing, ok := obj.(*networkingv1.Ingress)
	if !ok {
		return nil, fmt.Errorf("object type is not *networkingv1.Ingress")
	}
	return h.applyIngress(ing)
}

// ApplyFromUnstructured applies ingress from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*networkingv1.Ingress, error) {
	ing := &networkingv1.Ingress{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ing)
	if err != nil {
		return nil, err
	}
	return h.applyIngress(ing)
}

// ApplyFromMap applies ingress from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*networkingv1.Ingress, error) {
	ing := &networkingv1.Ingress{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ing)
	if err != nil {
		return nil, err
	}
	return h.applyIngress(ing)
}

// applyIngress
func (h *Handler) applyIngress(ing *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	_, err := h.createIngress(ing)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateIngress(ing)
	}
	return ing, err
}
