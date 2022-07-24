package ingressclass

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates ingressclass from type string, []byte, *networkingv1.IngressClass,
// networkingv1.IngressClass, runtime.Object or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*networkingv1.IngressClass, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *networkingv1.IngressClass:
		return h.UpdateFromObject(val)
	case networkingv1.IngressClass:
		return h.UpdateFromObject(&val)
	case runtime.Object:
		return h.UpdateFromObject(val)
	case map[string]interface{}:
		return h.UpdateFromUnstructured(val)
	default:
		return nil, ERR_TYPE_UPDATE
	}
}

// UpdateFromFile updates ingressclass from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*networkingv1.IngressClass, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates ingressclass from bytes.
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
	return h.updateIngressclass(ingc)
}

// UpdateFromObject updates ingressclass from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*networkingv1.IngressClass, error) {
	ingc, ok := obj.(*networkingv1.IngressClass)
	if !ok {
		return nil, fmt.Errorf("object is not *networkingv1.IngressClass")
	}
	return h.updateIngressclass(ingc)
}

// UpdateFromUnstructured updates ingressclass from map[string]interface{}.
func (h *Handler) UpdateFromUnstructured(u map[string]interface{}) (*networkingv1.IngressClass, error) {
	ingc := &networkingv1.IngressClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ingc)
	if err != nil {
		return nil, err
	}
	return h.updateIngressclass(ingc)
}

// updateIngressclass
func (h *Handler) updateIngressclass(ingc *networkingv1.IngressClass) (*networkingv1.IngressClass, error) {
	ingc.ResourceVersion = ""
	ingc.UID = ""
	return h.clientset.NetworkingV1().IngressClasses().Update(h.ctx, ingc, h.Options.UpdateOptions)
}
