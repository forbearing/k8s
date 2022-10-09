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

// Update updates ingress from type string, []byte, *networkingv1.Ingress,
// networkingv1.Ingress, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*networkingv1.Ingress, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *networkingv1.Ingress:
		return h.UpdateFromObject(val)
	case networkingv1.Ingress:
		return h.UpdateFromObject(&val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.UpdateFromObject(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates ingress from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*networkingv1.Ingress, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates ingress from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*networkingv1.Ingress, error) {
	ingJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ing := &networkingv1.Ingress{}
	if err = json.Unmarshal(ingJson, ing); err != nil {
		return nil, err
	}
	return h.updateIngress(ing)
}

// UpdateFromObject updates ingress from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*networkingv1.Ingress, error) {
	ing, ok := obj.(*networkingv1.Ingress)
	if !ok {
		return nil, fmt.Errorf("object type is not *networkingv1.Ingress")
	}
	return h.updateIngress(ing)
}

// UpdateFromUnstructured updates ingress from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*networkingv1.Ingress, error) {
	ing := &networkingv1.Ingress{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ing)
	if err != nil {
		return nil, err
	}
	return h.updateIngress(ing)
}

// UpdateFromMap updates ingress from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*networkingv1.Ingress, error) {
	ing := &networkingv1.Ingress{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ing)
	if err != nil {
		return nil, err
	}
	return h.updateIngress(ing)
}

// updateIngress
func (h *Handler) updateIngress(ing *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	namespace := ing.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	ing.ResourceVersion = ""
	ing.UID = ""
	return h.clientset.NetworkingV1().Ingresses(namespace).Update(h.ctx, ing, h.Options.UpdateOptions)
}
