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

// Delete deletes ingress from type string, []byte, *networkingv1.Ingress,
// networkingv1.Ingress, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a ingress from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *networkingv1.Ingress:
		return h.DeleteFromObject(val)
	case networkingv1.Ingress:
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

// DeleteByName deletes ingress by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.NetworkingV1().Ingresses(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes ingress from yaml or json file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes ingress from bytes data.
func (h *Handler) DeleteFromBytes(data []byte) error {
	ingJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	ing := &networkingv1.Ingress{}
	if err = json.Unmarshal(ingJson, ing); err != nil {
		return err
	}
	return h.deleteIngress(ing)
}

// DeleteFromObject deletes ingress from metav1.Object or runtime.Object.
func (h *Handler) DeleteFromObject(obj interface{}) error {
	ing, ok := obj.(*networkingv1.Ingress)
	if !ok {
		return fmt.Errorf("object type is not *networkingv1.Ingress")
	}
	return h.deleteIngress(ing)
}

// DeleteFromUnstructured deletes ingress from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	ing := &networkingv1.Ingress{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ing)
	if err != nil {
		return err
	}
	return h.deleteIngress(ing)
}

// DeleteFromMap deletes ingress from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	ing := &networkingv1.Ingress{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ing)
	if err != nil {
		return err
	}
	return h.deleteIngress(ing)
}

// deleteIngress
func (h *Handler) deleteIngress(ing *networkingv1.Ingress) error {
	namespace := ing.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.NetworkingV1().Ingresses(namespace).Delete(h.ctx, ing.Name, h.Options.DeleteOptions)
}
