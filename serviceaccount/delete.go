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

// Delete deletes serviceaccount from type string, []byte, *corev1.ServiceAccount,
// corev1.ServiceAccount, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a serviceaccount from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *corev1.ServiceAccount:
		return h.DeleteFromObject(val)
	case corev1.ServiceAccount:
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

// DeleteByName deletes serviceaccount by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().ServiceAccounts(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes serviceaccount from yaml or json file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes serviceaccount from bytes data.
func (h *Handler) DeleteFromBytes(data []byte) error {
	saJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	sa := &corev1.ServiceAccount{}
	if err = json.Unmarshal(saJson, sa); err != nil {
		return err
	}
	return h.deleteSA(sa)
}

// DeleteFromObject deletes serviceaccount from metav1.Object or runtime.Object.
func (h *Handler) DeleteFromObject(obj interface{}) error {
	sa, ok := obj.(*corev1.ServiceAccount)
	if !ok {
		return fmt.Errorf("object type is not *corev1.ServiceAccount")
	}
	return h.deleteSA(sa)
}

// DeleteFromUnstructured deletes serviceaccount from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	sa := &corev1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sa)
	if err != nil {
		return err
	}
	return h.deleteSA(sa)
}

// DeleteFromMap deletes serviceaccount from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	sa := &corev1.ServiceAccount{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sa)
	if err != nil {
		return err
	}
	return h.deleteSA(sa)
}

// deleteSA
func (h *Handler) deleteSA(sa *corev1.ServiceAccount) error {
	namespace := sa.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().ServiceAccounts(namespace).Delete(h.ctx, sa.Name, h.Options.DeleteOptions)
}
