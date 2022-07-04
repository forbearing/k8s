package ingress

import (
	"encoding/json"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create ingress from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*networkingv1.Ingress, error) {
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

	return h.clientset.NetworkingV1().Ingresses(namespace).Create(h.ctx, ingress, h.Options.CreateOptions)
}

// CreateFromBytes create ingress from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*networkingv1.Ingress, error) {
	ingressJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ingress := &networkingv1.Ingress{}
	err = json.Unmarshal(ingressJson, ingress)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(ingress.Namespace) != 0 {
		namespace = ingress.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.NetworkingV1().Ingresses(namespace).Create(h.ctx, ingress, h.Options.CreateOptions)
}

// CreateFromFile create ingress from yaml file.
func (h *Handler) CreateFromFile(filename string) (*networkingv1.Ingress, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create ingress from file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*networkingv1.Ingress, error) {
	return h.CreateFromFile(filename)
}
