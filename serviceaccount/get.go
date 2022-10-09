package serviceaccount

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

// Get gets serviceaccount from type string, []byte, *corev1.ServiceAccount,
// corev1.ServiceAccount, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a serviceaccount from file path.
func (h *Handler) Get(obj interface{}) (*corev1.ServiceAccount, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *corev1.ServiceAccount:
		return h.GetFromObject(val)
	case corev1.ServiceAccount:
		return h.GetFromObject(&val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	case metav1.Object, runtime.Object:
		return h.GetFromObject(val)
	default:
		return nil, ErrInvalidGetType
	}
}

// GetByName gets serviceaccount by name.
func (h *Handler) GetByName(name string) (*corev1.ServiceAccount, error) {
	return h.clientset.CoreV1().ServiceAccounts(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets serviceaccount from yaml or json file.
func (h *Handler) GetFromFile(filename string) (*corev1.ServiceAccount, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets serviceaccount from bytes data.
func (h *Handler) GetFromBytes(data []byte) (*corev1.ServiceAccount, error) {
	saJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sa := &corev1.ServiceAccount{}
	if err = json.Unmarshal(saJson, sa); err != nil {
		return nil, err
	}
	return h.getSA(sa)
}

// GetFromObject gets serviceaccount from metav1.Object or runtime.Object.
func (h *Handler) GetFromObject(obj interface{}) (*corev1.ServiceAccount, error) {
	sa, ok := obj.(*corev1.ServiceAccount)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.ServiceAccount")
	}
	return h.getSA(sa)
}

// GetFromUnstructured gets serviceaccount from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*corev1.ServiceAccount, error) {
	sa := &corev1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sa)
	if err != nil {
		return nil, err
	}
	return h.getSA(sa)
}

// GetFromMap gets serviceaccount from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*corev1.ServiceAccount, error) {
	sa := &corev1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sa)
	if err != nil {
		return nil, err
	}
	return h.getSA(sa)
}

// getSA
// It's necessary to get a new serviceaccount resource from a old serviceaccount resource,
// because old serviceaccount usually don't have serviceaccount.Status field.
func (h *Handler) getSA(sa *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
	namespace := sa.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().ServiceAccounts(namespace).Get(h.ctx, sa.Name, h.Options.GetOptions)
}
