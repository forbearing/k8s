package ingress

import (
	"encoding/json"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get ingress from bytes.
func (h *Handler) GetFromBytes(data []byte) (*networkingv1.Ingress, error) {
	ingressJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}
	ingress := &networkingv1.Ingress{}
	if err = json.Unmarshal(ingressJson, ingress); err != nil {
		return nil, err
	}

	var namespace string
	if len(ingress.Namespace) != 0 {
		namespace = ingress.Namespace
	} else {
		namespace = h.namespace
	}
	return h.WithNamespace(namespace).GetByName(ingress.Name)
}

// GetFromFile get ingress from yaml file.
func (h *Handler) GetFromFile(filename string) (*networkingv1.Ingress, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get ingress by name.
func (h *Handler) GetByName(name string) (*networkingv1.Ingress, error) {
	return h.clientset.NetworkingV1().Ingresses(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get ingress by name, alias to "GetByName".
func (h *Handler) Get(name string) (*networkingv1.Ingress, error) {
	return h.GetByName(name)
}
