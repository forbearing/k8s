package ingress

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets ingress from type string, []byte, *networkingv1.Ingress,
// networkingv1.Ingress, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a ingress from file path.
func (h *Handler) Get(obj interface{}) (*networkingv1.Ingress, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *networkingv1.Ingress:
		return h.GetFromObject(val)
	case networkingv1.Ingress:
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

// GetByName gets ingress by name.
func (h *Handler) GetByName(name string) (*networkingv1.Ingress, error) {
	return h.clientset.NetworkingV1().Ingresses(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets ingress from yaml or json file.
func (h *Handler) GetFromFile(filename string) (*networkingv1.Ingress, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets ingress from bytes data.
func (h *Handler) GetFromBytes(data []byte) (*networkingv1.Ingress, error) {
	ingJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ing := &networkingv1.Ingress{}
	if err = json.Unmarshal(ingJson, ing); err != nil {
		return nil, err
	}
	return h.getIngress(ing)
}

// GetFromObject gets ingress from metav1.Object or runtime.Object.
func (h *Handler) GetFromObject(obj interface{}) (*networkingv1.Ingress, error) {
	ing, ok := obj.(*networkingv1.Ingress)
	if !ok {
		return nil, fmt.Errorf("object type is not *networkingv1.Ingress")
	}
	return h.getIngress(ing)
}

// GetFromUnstructured gets ingress from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*networkingv1.Ingress, error) {
	ing := &networkingv1.Ingress{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ing)
	if err != nil {
		return nil, err
	}
	return h.getIngress(ing)
}

// GetFromMap gets ingress from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*networkingv1.Ingress, error) {
	ing := &networkingv1.Ingress{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ing)
	if err != nil {
		return nil, err
	}
	return h.getIngress(ing)
}

// getIngress
// It's necessary to get a new ingress resource from a old ingress resource,
// because old ingress usually don't have ingress.Status field.
func (h *Handler) getIngress(ing *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	namespace := ing.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.NetworkingV1().Ingresses(namespace).Get(h.ctx, ing.Name, h.Options.GetOptions)
}
