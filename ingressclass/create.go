package ingressclass

import (
	"encoding/json"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create ingressclass from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*networkingv1.IngressClass, error) {
	ingressclass := &networkingv1.IngressClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, ingressclass)
	if err != nil {
		return nil, err
	}

	return h.clientset.NetworkingV1().IngressClasses().Create(h.ctx, ingressclass, h.Options.CreateOptions)
}

// CreateFromBytes create ingressclass from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*networkingv1.IngressClass, error) {
	ingcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ingc := &networkingv1.IngressClass{}
	err = json.Unmarshal(ingcJson, ingc)
	if err != nil {
		return nil, err
	}

	return h.clientset.NetworkingV1().IngressClasses().Create(h.ctx, ingc, h.Options.CreateOptions)
}

// CreateFromFile create ingressclass from yaml file.
func (h *Handler) CreateFromFile(filename string) (*networkingv1.IngressClass, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create ingressclass from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*networkingv1.IngressClass, error) {
	return h.CreateFromFile(filename)
}
