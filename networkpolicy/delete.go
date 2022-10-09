package networkpolicy

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

// Delete deletes networkpolicy from type string, []byte, *networkingv1.NetworkPolicy,
// networkingv1.NetworkPolicy, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a networkpolicy from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *networkingv1.NetworkPolicy:
		return h.DeleteFromObject(val)
	case networkingv1.NetworkPolicy:
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

// DeleteByName deletes networkpolicy by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.NetworkingV1().NetworkPolicies(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes networkpolicy from yaml or json file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes networkpolicy from bytes data.
func (h *Handler) DeleteFromBytes(data []byte) error {
	netpolJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	netpol := &networkingv1.NetworkPolicy{}
	if err = json.Unmarshal(netpolJson, netpol); err != nil {
		return err
	}
	return h.deleteNetpol(netpol)
}

// DeleteFromObject deletes networkpolicy from metav1.Object or runtime.Object.
func (h *Handler) DeleteFromObject(obj interface{}) error {
	netpol, ok := obj.(*networkingv1.NetworkPolicy)
	if !ok {
		return fmt.Errorf("object type is not *networkingv1.NetworkPolicy")
	}
	return h.deleteNetpol(netpol)
}

// DeleteFromUnstructured deletes networkpolicy from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	netpol := &networkingv1.NetworkPolicy{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), netpol)
	if err != nil {
		return err
	}
	return h.deleteNetpol(netpol)
}

// DeleteFromMap deletes networkpolicy from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	netpol := &networkingv1.NetworkPolicy{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, netpol)
	if err != nil {
		return err
	}
	return h.deleteNetpol(netpol)
}

// deleteNetpol
func (h *Handler) deleteNetpol(netpol *networkingv1.NetworkPolicy) error {
	namespace := netpol.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.NetworkingV1().NetworkPolicies(namespace).Delete(h.ctx, netpol.Name, h.Options.DeleteOptions)
}
