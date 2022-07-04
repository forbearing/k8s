package ingressclass

import (
	"encoding/json"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete ingressclass from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	ingcJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	ingc := &networkingv1.IngressClass{}
	err = json.Unmarshal(ingcJson, ingc)
	if err != nil {
		return err
	}

	return h.DeleteByName(ingc.Name)
}

// DeleteFromFile delete ingressclass from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete ingressclass by name
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.NetworkingV1().IngressClasses().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete ingressclass by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
