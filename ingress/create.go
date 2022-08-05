package ingress

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

// Create creates ingress from type string, []byte, *networkingv1.Ingress,
// networkingv1.Ingress, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*networkingv1.Ingress, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *networkingv1.Ingress:
		return h.CreateFromObject(val)
	case networkingv1.Ingress:
		return h.CreateFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.CreateFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.CreateFromObject(val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	default:
		return nil, ErrInvalidCreateType
	}
}

// CreateFromFile creates ingress from yaml file.
func (h *Handler) CreateFromFile(filename string) (*networkingv1.Ingress, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates ingress from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*networkingv1.Ingress, error) {
	ingJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ing := &networkingv1.Ingress{}
	if err = json.Unmarshal(ingJson, ing); err != nil {
		return nil, err
	}
	return h.createIngress(ing)
}

// CreateFromObject creates ingress from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*networkingv1.Ingress, error) {
	ing, ok := obj.(*networkingv1.Ingress)
	if !ok {
		return nil, fmt.Errorf("object type is not *networkingv1.Ingress")
	}
	return h.createIngress(ing)
}

// CreateFromUnstructured creates ingress from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*networkingv1.Ingress, error) {
	ing := &networkingv1.Ingress{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ing)
	if err != nil {
		return nil, err
	}
	return h.createIngress(ing)
}

// CreateFromMap creates ingress from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*networkingv1.Ingress, error) {
	ing := &networkingv1.Ingress{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ing)
	if err != nil {
		return nil, err
	}
	return h.createIngress(ing)
}

// createIngress
func (h *Handler) createIngress(ing *networkingv1.Ingress) (*networkingv1.Ingress, error) {
	var namespace string
	if len(ing.Namespace) != 0 {
		namespace = ing.Namespace
	} else {
		namespace = h.namespace
	}
	ing.ResourceVersion = ""
	ing.UID = ""
	return h.clientset.NetworkingV1().Ingresses(namespace).Create(h.ctx, ing, h.Options.CreateOptions)
}
