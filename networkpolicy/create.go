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

// Create creates networkpolicy from type string, []byte, *networkingv1.NetworkPolicy,
// networkingv1.NetworkPolicy, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*networkingv1.NetworkPolicy, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *networkingv1.NetworkPolicy:
		return h.CreateFromObject(val)
	case networkingv1.NetworkPolicy:
		return h.CreateFromObject(&val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates networkpolicy from yaml file.
func (h *Handler) CreateFromFile(filename string) (*networkingv1.NetworkPolicy, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates networkpolicy from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*networkingv1.NetworkPolicy, error) {
	netpolJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	netpol := &networkingv1.NetworkPolicy{}
	err = json.Unmarshal(netpolJson, netpol)
	if err != nil {
		return nil, err
	}
	return h.createNetpol(netpol)
}

// CreateFromObject creates networkpolicy from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*networkingv1.NetworkPolicy, error) {
	netpol, ok := obj.(*networkingv1.NetworkPolicy)
	if !ok {
		return nil, fmt.Errorf("object type is not *networkingv1.NetworkPolicy")
	}
	return h.createNetpol(netpol)
}

// CreateFromUnstructured creates networkpolicy from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*networkingv1.NetworkPolicy, error) {
	netpol := &networkingv1.NetworkPolicy{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), netpol)
	if err != nil {
		return nil, err
	}
	return h.createNetpol(netpol)
}

// CreateFromMap creates networkpolicy from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*networkingv1.NetworkPolicy, error) {
	netpol := &networkingv1.NetworkPolicy{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, netpol)
	if err != nil {
		return nil, err
	}
	return h.createNetpol(netpol)
}

// createNetpol
func (h *Handler) createNetpol(netpol *networkingv1.NetworkPolicy) (*networkingv1.NetworkPolicy, error) {
	var namespace string
	if len(netpol.Namespace) != 0 {
		namespace = netpol.Namespace
	} else {
		namespace = h.namespace
	}
	netpol.ResourceVersion = ""
	netpol.UID = ""
	return h.clientset.NetworkingV1().NetworkPolicies(namespace).Create(h.ctx, netpol, h.Options.CreateOptions)
}
