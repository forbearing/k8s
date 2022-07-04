package networkpolicy

import (
	"encoding/json"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create networkpolicy from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*networkingv1.NetworkPolicy, error) {
	netpol := &networkingv1.NetworkPolicy{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, netpol)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(netpol.Namespace) != 0 {
		namespace = netpol.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.NetworkingV1().NetworkPolicies(namespace).Create(h.ctx, netpol, h.Options.CreateOptions)
}

// CreateFromBytes create networkpolicy from bytes.
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

	var namespace string
	if len(netpol.Namespace) != 0 {
		namespace = netpol.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.NetworkingV1().NetworkPolicies(namespace).Create(h.ctx, netpol, h.Options.CreateOptions)
}

// CreateFromFile create networkpolicy from yaml file.
func (h *Handler) CreateFromFile(filename string) (*networkingv1.NetworkPolicy, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create networkpolicy from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*networkingv1.NetworkPolicy, error) {
	return h.CreateFromFile(filename)
}
