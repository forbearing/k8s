package networkpolicy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Delete deletes networkpolicy from type string, []byte, *networkingv1.NetworkPolicy,
// networkingv1.NetworkPolicy, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

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

// DeleteByName deletes networkpolicy by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.NetworkingV1().NetworkPolicies(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes networkpolicy from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes networkpolicy from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	netpolJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	netpol := &networkingv1.NetworkPolicy{}
	err = json.Unmarshal(netpolJson, netpol)
	if err != nil {
		return err
	}
	return h.deleteNetpol(netpol)
}

// DeleteFromObject deletes networkpolicy from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
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
	var namespace string
	if len(netpol.Namespace) != 0 {
		namespace = netpol.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.NetworkingV1().NetworkPolicies(namespace).Delete(h.ctx, netpol.Name, h.Options.DeleteOptions)
}
