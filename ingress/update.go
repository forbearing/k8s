package ingress

import (
	"encoding/json"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update ingress from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*networkingv1.Ingress, error) {
	ingress := &networkingv1.Ingress{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, ingress)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(ingress.Namespace) != 0 {
		namespace = ingress.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.NetworkingV1().Ingresses(namespace).Update(h.ctx, ingress, h.Options.UpdateOptions)
}

// UpdateFromBytes update ingress from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*networkingv1.Ingress, error) {
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

	return h.clientset.NetworkingV1().Ingresses(namespace).Update(h.ctx, ingress, h.Options.UpdateOptions)
}

// UpdateFromFile update ingress from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*networkingv1.Ingress, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update ingress from file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*networkingv1.Ingress, error) {
	return h.UpdateFromFile(filename)
}
