package ingressclass

import (
	"encoding/json"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update ingressclass from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*networkingv1.IngressClass, error) {
	ingressclass := &networkingv1.IngressClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, ingressclass)
	if err != nil {
		return nil, err
	}

	return h.clientset.NetworkingV1().IngressClasses().Update(h.ctx, ingressclass, h.Options.UpdateOptions)
}

// UpdateFromBytes update ingressclass from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*networkingv1.IngressClass, error) {
	ingcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ingc := &networkingv1.IngressClass{}
	err = json.Unmarshal(ingcJson, ingc)
	if err != nil {
		return nil, err
	}

	return h.clientset.NetworkingV1().IngressClasses().Update(h.ctx, ingc, h.Options.UpdateOptions)
}

// UpdateFromFile update ingressclass from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*networkingv1.IngressClass, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update ingressclass from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*networkingv1.IngressClass, error) {
	return h.UpdateFromFile(filename)
}
