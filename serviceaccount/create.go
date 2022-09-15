package serviceaccount

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates serviceaccount from type string, []byte, *corev1.ServiceAccount,
// corev1.ServiceAccount, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*corev1.ServiceAccount, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *corev1.ServiceAccount:
		return h.CreateFromObject(val)
	case corev1.ServiceAccount:
		return h.CreateFromObject(&val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	case runtime.Object:
		return h.CreateFromObject(val)
	default:
		return nil, ErrInvalidCreateType
	}
}

// CreateFromFile creates serviceaccount from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.ServiceAccount, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates serviceaccount from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.ServiceAccount, error) {
	saJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sa := &corev1.ServiceAccount{}
	if err = json.Unmarshal(saJson, sa); err != nil {
		return nil, err
	}
	return h.createSA(sa)
}

// CreateFromObject creates serviceaccount from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*corev1.ServiceAccount, error) {
	sa, ok := obj.(*corev1.ServiceAccount)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.ServiceAccount")
	}
	return h.createSA(sa)
}

// CreateFromUnstructured creates serviceaccount from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*corev1.ServiceAccount, error) {
	sa := &corev1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sa)
	if err != nil {
		return nil, err
	}
	return h.createSA(sa)
}

// CreateFromMap creates serviceaccount from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*corev1.ServiceAccount, error) {
	sa := &corev1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sa)
	if err != nil {
		return nil, err
	}
	return h.createSA(sa)
}

// createSA
func (h *Handler) createSA(sa *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
	var namespace string
	if len(sa.Namespace) != 0 {
		namespace = sa.Namespace
	} else {
		namespace = h.namespace
	}
	sa.ResourceVersion = ""
	sa.UID = ""
	return h.clientset.CoreV1().ServiceAccounts(namespace).Create(h.ctx, sa, h.Options.CreateOptions)
}
