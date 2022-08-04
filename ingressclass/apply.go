package ingressclass

import (
	"fmt"
	"reflect"

	networkingv1 "k8s.io/api/networking/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies ingressclass from type string, []byte, *networkingv1.IngressClass,
// networkingv1.IngressClass, runtime.Object, *unstructured.Unstructured,
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
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.ApplyFromUnstructured(val.(*unstructured.Unstructured))
		}
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

// ApplyFromFile applies ingressclass from yaml file.
func (h *Handler) ApplyFromFile(filename string) (ingc *networkingv1.IngressClass, err error) {
	ingc, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if ingressclass already exist, update it.
		ingc, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply ingressclass from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (ingc *networkingv1.IngressClass, err error) {
	ingc, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		ingc, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies ingressclass from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*networkingv1.IngressClass, error) {
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
