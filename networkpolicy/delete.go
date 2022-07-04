package networkpolicy

import (
	"encoding/json"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete networkpolicy from bytes.
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

	var namespace string
	if len(netpol.Namespace) != 0 {
		namespace = netpol.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(netpol.Name)
}

// DeleteFromFile delete networkpolicy from yaml file.
func (h *Handler) DeleteFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete networkpolicy by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.NetworkingV1().NetworkPolicies(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete networkpolicy by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
