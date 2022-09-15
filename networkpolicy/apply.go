package networkpolicy

import (
	"fmt"

	networkingv1 "k8s.io/api/networking/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies networkpolicy from type string, []byte, *networkingv1.NetworkPolicy,
// networkingv1.NetworkPolicy, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*networkingv1.NetworkPolicy, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *networkingv1.NetworkPolicy:
		return h.ApplyFromObject(val)
	case networkingv1.NetworkPolicy:
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

// ApplyFromFile applies networkpolicy from yaml file.
func (h *Handler) ApplyFromFile(filename string) (netpol *networkingv1.NetworkPolicy, err error) {
	netpol, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if networkpolicy already exist, update it.
		netpol, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply networkpolicy from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (netpol *networkingv1.NetworkPolicy, err error) {
	netpol, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		netpol, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies networkpolicy from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*networkingv1.NetworkPolicy, error) {
	netpol, ok := obj.(*networkingv1.NetworkPolicy)
	if !ok {
		return nil, fmt.Errorf("object type is not *networkingv1.NetworkPolicy")
	}
	return h.applyNetpol(netpol)
}

// ApplyFromUnstructured applies networkpolicy from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*networkingv1.NetworkPolicy, error) {
	netpol := &networkingv1.NetworkPolicy{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), netpol)
	if err != nil {
		return nil, err
	}
	return h.applyNetpol(netpol)
}

// ApplyFromMap applies networkpolicy from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*networkingv1.NetworkPolicy, error) {
	netpol := &networkingv1.NetworkPolicy{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, netpol)
	if err != nil {
		return nil, err
	}
	return h.applyNetpol(netpol)
}

// applyNetpol
func (h *Handler) applyNetpol(netpol *networkingv1.NetworkPolicy) (*networkingv1.NetworkPolicy, error) {
	_, err := h.createNetpol(netpol)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateNetpol(netpol)
	}
	return netpol, err
}
