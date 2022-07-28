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

// Update updates networkpolicy from type string, []byte, *networkingv1.NetworkPolicy,
// networkingv1.NetworkPolicy, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*networkingv1.NetworkPolicy, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *networkingv1.NetworkPolicy:
		return h.UpdateFromObject(val)
	case networkingv1.NetworkPolicy:
		return h.UpdateFromObject(&val)
	case runtime.Object:
		return h.UpdateFromObject(val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	default:
		return nil, ERR_TYPE_UPDATE
	}
}

// UpdateFromFile updates networkpolicy from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*networkingv1.NetworkPolicy, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates networkpolicy from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*networkingv1.NetworkPolicy, error) {
	netpolJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	netpol := &networkingv1.NetworkPolicy{}
	if err = json.Unmarshal(netpolJson, netpol); err != nil {
		return nil, err
	}
	return h.updateNetpol(netpol)
}

// UpdateFromObject updates networkpolicy from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*networkingv1.NetworkPolicy, error) {
	netpol, ok := obj.(*networkingv1.NetworkPolicy)
	if !ok {
		return nil, fmt.Errorf("object type is not *networkingv1.NetworkPolicy")
	}
	return h.updateNetpol(netpol)
}

// UpdateFromUnstructured updates networkpolicy from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*networkingv1.NetworkPolicy, error) {
	netpol := &networkingv1.NetworkPolicy{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), netpol)
	if err != nil {
		return nil, err
	}
	return h.updateNetpol(netpol)
}

// UpdateFromMap updates networkpolicy from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*networkingv1.NetworkPolicy, error) {
	netpol := &networkingv1.NetworkPolicy{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, netpol)
	if err != nil {
		return nil, err
	}
	return h.updateNetpol(netpol)
}

// updateNetpol
func (h *Handler) updateNetpol(netpol *networkingv1.NetworkPolicy) (*networkingv1.NetworkPolicy, error) {
	var namespace string
	if len(netpol.Namespace) != 0 {
		namespace = netpol.Namespace
	} else {
		namespace = h.namespace
	}
	netpol.ResourceVersion = ""
	netpol.UID = ""
	return h.clientset.NetworkingV1().NetworkPolicies(namespace).Update(h.ctx, netpol, h.Options.UpdateOptions)
}
