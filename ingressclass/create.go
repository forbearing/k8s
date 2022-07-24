package ingressclass

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates ingressclass from type string, []byte, *networkingv1.IngressClass,
// networkingv1.IngressClass, runtime.Object or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*networkingv1.IngressClass, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *networkingv1.IngressClass:
		return h.CreateFromObject(val)
	case networkingv1.IngressClass:
		return h.CreateFromObject(&val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case map[string]interface{}:
		return h.CreateFromUnstructured(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates ingressclass from yaml file.
func (h *Handler) CreateFromFile(filename string) (*networkingv1.IngressClass, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates ingressclass from bytes.
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
	return h.createIngressclass(ingc)
}

// CreateFromObject creates ingressclass from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*networkingv1.IngressClass, error) {
	ingc, ok := obj.(*networkingv1.IngressClass)
	if !ok {
		return nil, fmt.Errorf("object is not *networkingv1.IngressClass")
	}
	return h.createIngressclass(ingc)
}

// CreateFromUnstructured creates ingressclass from map[string]interface{}.
func (h *Handler) CreateFromUnstructured(u map[string]interface{}) (*networkingv1.IngressClass, error) {
	ingc := &networkingv1.IngressClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ingc)
	if err != nil {
		return nil, err
	}
	return h.createIngressclass(ingc)
}

// createIngressclass
func (h *Handler) createIngressclass(ingc *networkingv1.IngressClass) (*networkingv1.IngressClass, error) {
	ingc.ResourceVersion = ""
	ingc.UID = ""
	return h.clientset.NetworkingV1().IngressClasses().Create(h.ctx, ingc, h.Options.CreateOptions)
}
