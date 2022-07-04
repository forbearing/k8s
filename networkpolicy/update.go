package networkpolicy

import (
	"encoding/json"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update networkpolicy from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*networkingv1.NetworkPolicy, error) {
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

	return h.clientset.NetworkingV1().NetworkPolicies(namespace).Update(h.ctx, netpol, h.Options.UpdateOptions)
}

// UpdateFromBytes update networkpolicy from bytesA.
func (h *Handler) UpdateFromBytes(data []byte) (*networkingv1.NetworkPolicy, error) {
	netpolJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	netpol := &networkingv1.NetworkPolicy{}
	if err = json.Unmarshal(netpolJson, netpol); err != nil {
		return nil, err
	}

	var namespace string
	if len(netpol.Namespace) != 0 {
		namespace = netpol.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.NetworkingV1().NetworkPolicies(namespace).Update(h.ctx, netpol, h.Options.UpdateOptions)
}

// UpdateFromFile update networkpolicy from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*networkingv1.NetworkPolicy, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update networkpolicy from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*networkingv1.NetworkPolicy, error) {
	return h.UpdateFromFile(filename)
}
