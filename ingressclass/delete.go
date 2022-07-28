package ingressclass

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Delete deletes ingressclass from type string, []byte, *networkingv1.IngressClass,
// networkingv1.IngressClass, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a ingressclass from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *networkingv1.IngressClass:
		return h.DeleteFromObject(val)
	case networkingv1.IngressClass:
		return h.DeleteFromObject(&val)
	case runtime.Object:
		return h.DeleteFromObject(val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	default:
		return ERR_TYPE_DELETE
	}
}

// DeleteByName deletes ingressclass by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.NetworkingV1().IngressClasses().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes ingressclass from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes ingressclass from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	ingcJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	ingc := &networkingv1.IngressClass{}
	err = json.Unmarshal(ingcJson, ingc)
	if err != nil {
		return err
	}
	return h.deleteIngressclass(ingc)
}

// DeleteFromObject deletes ingressclass from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	ingc, ok := obj.(*networkingv1.IngressClass)
	if !ok {
		return fmt.Errorf("object type is not *networkingv1.IngressClass")
	}
	return h.deleteIngressclass(ingc)
}

// DeleteFromUnstructured deletes ingressclass from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	ingc := &networkingv1.IngressClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ingc)
	if err != nil {
		return err
	}
	return h.deleteIngressclass(ingc)
}

// DeleteFromMap deletes ingressclass from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	ingc := &networkingv1.IngressClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ingc)
	if err != nil {
		return err
	}
	return h.deleteIngressclass(ingc)
}

// deleteIngressclass
func (h *Handler) deleteIngressclass(ingc *networkingv1.IngressClass) error {
	return h.clientset.NetworkingV1().IngressClasses().Delete(h.ctx, ingc.Name, h.Options.DeleteOptions)
}
