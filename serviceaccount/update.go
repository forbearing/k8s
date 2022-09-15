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

// Update updates serviceaccount from type string, []byte, *corev1.ServiceAccount,
// corev1.ServiceAccount, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*corev1.ServiceAccount, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *corev1.ServiceAccount:
		return h.UpdateFromObject(val)
	case corev1.ServiceAccount:
		return h.UpdateFromObject(&val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	case runtime.Object:
		return h.UpdateFromObject(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates serviceaccount from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.ServiceAccount, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates serviceaccount from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.ServiceAccount, error) {
	saJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sa := &corev1.ServiceAccount{}
	if err = json.Unmarshal(saJson, sa); err != nil {
		return nil, err
	}
	return h.updateSA(sa)
}

// UpdateFromObject updates serviceaccount from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*corev1.ServiceAccount, error) {
	sa, ok := obj.(*corev1.ServiceAccount)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.ServiceAccount")
	}
	return h.updateSA(sa)
}

// UpdateFromUnstructured updates serviceaccount from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*corev1.ServiceAccount, error) {
	sa := &corev1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sa)
	if err != nil {
		return nil, err
	}
	return h.updateSA(sa)
}

// UpdateFromMap updates serviceaccount from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*corev1.ServiceAccount, error) {
	sa := &corev1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sa)
	if err != nil {
		return nil, err
	}
	return h.updateSA(sa)
}

// updateSA
func (h *Handler) updateSA(sa *corev1.ServiceAccount) (*corev1.ServiceAccount, error) {
	var namespace string
	if len(sa.Namespace) != 0 {
		namespace = sa.Namespace
	} else {
		namespace = h.namespace
	}
	sa.ResourceVersion = ""
	sa.UID = ""
	return h.clientset.CoreV1().ServiceAccounts(namespace).Update(h.ctx, sa, h.Options.UpdateOptions)
}
