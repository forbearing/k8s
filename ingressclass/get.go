package ingressclass

import (
	"encoding/json"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get ingressclass from bytes.
func (h *Handler) GetFromBytes(data []byte) (*networkingv1.IngressClass, error) {
	ingcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ingc := &networkingv1.IngressClass{}
	err = json.Unmarshal(ingcJson, ingc)
	if err != nil {
		return nil, err
	}

	return h.GetByName(ingc.Name)
}

// GetFromFile get ingressclass from yaml file.
func (h *Handler) GetFromFile(filename string) (*networkingv1.IngressClass, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get ingressclass by name
func (h *Handler) GetByName(name string) (*networkingv1.IngressClass, error) {
	return h.clientset.NetworkingV1().IngressClasses().Get(h.ctx, name, h.Options.GetOptions)
}

// Get get ingressclass by name, alias to "GetByName".
func (h *Handler) Get(name string) (*networkingv1.IngressClass, error) {
	return h.clientset.NetworkingV1().IngressClasses().Get(h.ctx, name, h.Options.GetOptions)
}
