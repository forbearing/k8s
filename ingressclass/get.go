package ingressclass

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets ingressclass from type string, []byte, *networkingv1.IngressClass,
// networkingv1.IngressClass, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a ingressclass from file path.
func (h *Handler) Get(obj interface{}) (*networkingv1.IngressClass, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *networkingv1.IngressClass:
		return h.GetFromObject(val)
	case networkingv1.IngressClass:
		return h.GetFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.GetFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.GetFromObject(val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets ingressclass by name.
func (h *Handler) GetByName(name string) (*networkingv1.IngressClass, error) {
	return h.clientset.NetworkingV1().IngressClasses().Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets ingressclass from yaml file.
func (h *Handler) GetFromFile(filename string) (*networkingv1.IngressClass, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets ingressclass from bytes.
func (h *Handler) GetFromBytes(data []byte) (*networkingv1.IngressClass, error) {
	ingcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ingc := &networkingv1.IngressClass{}
	if err = json.Unmarshal(ingcJson, ingc); err != nil {
		return nil, err
	}
	return h.getIngressclass(ingc)
}

// GetFromObject gets ingressclass from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*networkingv1.IngressClass, error) {
	ingc, ok := obj.(*networkingv1.IngressClass)
	if !ok {
		return nil, fmt.Errorf("object type is not *networkingv1.IngressClass")
	}
	return h.getIngressclass(ingc)
}

// GetFromUnstructured gets ingressclass from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*networkingv1.IngressClass, error) {
	ingc := &networkingv1.IngressClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ingc)
	if err != nil {
		return nil, err
	}
	return h.getIngressclass(ingc)
}

// GetFromMap gets ingressclass from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*networkingv1.IngressClass, error) {
	ingc := &networkingv1.IngressClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ingc)
	if err != nil {
		return nil, err
	}
	return h.getIngressclass(ingc)
}

// getIngressclass
// It's necessary to get a new ingressclass resource from a old ingressclass resource,
// because old ingressclass usually don't have ingressclass.Status field.
func (h *Handler) getIngressclass(ingc *networkingv1.IngressClass) (*networkingv1.IngressClass, error) {
	return h.clientset.NetworkingV1().IngressClasses().Get(h.ctx, ingc.Name, h.Options.GetOptions)
}
