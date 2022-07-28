package serviceaccount

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies serviceaccount from type string, []byte, *corev1.ServiceAccount,
// corev1.ServiceAccount, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*corev1.ServiceAccount, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *corev1.ServiceAccount:
		return h.ApplyFromObject(val)
	case corev1.ServiceAccount:
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

// ApplyFromFile applies serviceaccount from yaml file.
func (h *Handler) ApplyFromFile(filename string) (sa *corev1.ServiceAccount, err error) {
	sa, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if serviceaccount already exist, update it.
		sa, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply serviceaccount from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (sa *corev1.ServiceAccount, err error) {
	sa, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		sa, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies serviceaccount from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*corev1.ServiceAccount, error) {
	sa, ok := obj.(*corev1.ServiceAccount)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.ServiceAccount")
	}
	return h.applySA(sa)
}

// ApplyFromUnstructured applies serviceaccount from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*corev1.ServiceAccount, error) {
	sa := &corev1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sa)
	if err != nil {
		return nil, err
	}
	return h.applySA(sa)
}

// ApplyFromMap applies serviceaccount from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*corev1.ServiceAccount, error) {
	sa := &corev1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sa)
	if err != nil {
		return nil, err
	}
	return h.applySA(sa)
}

// applySA
func (h *Handler) applySA(sa *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
	_, err := h.createSA(sa)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateSA(sa)
	}
	return sa, err
}
