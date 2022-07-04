package networkpolicy

import (
	"encoding/json"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get networkpolicy from bytes.
func (h *Handler) GetFromBytes(data []byte) (*networkingv1.NetworkPolicy, error) {
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

	return h.WithNamespace(namespace).GetByName(netpol.Name)
}

// GetFromBytes get networkpolicy from yaml file.
func (h *Handler) GetFromFile(filename string) (*networkingv1.NetworkPolicy, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get networkpolicy by name.
func (h *Handler) GetByName(name string) (*networkingv1.NetworkPolicy, error) {
	return h.clientset.NetworkingV1().NetworkPolicies(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get networkpolicy by name, alias to "GetByName".
func (h *Handler) Get(name string) (*networkingv1.NetworkPolicy, error) {
	return h.GetByName(name)
}
