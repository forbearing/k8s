package ingress

import (
	"encoding/json"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete ingress from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	ingressJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	ingress := &networkingv1.Ingress{}
	err = json.Unmarshal(ingressJson, ingress)
	if err != nil {
		return err
	}

	var namespace string
	if len(ingress.Namespace) != 0 {
		namespace = ingress.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(ingress.Name)
}

// DeleteFromFile delete ingress from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete ingress by name
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.NetworkingV1().Ingresses(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete ingress by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
