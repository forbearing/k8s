package dynamic

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets unstructured k8s resource from type string, []byte,
// runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.

// If psssed parameter type is string, it will call GetByName insteard of GetFromFile.
// You  should always explicitly call GetFromFile to delete a unstructured object
// from filename.
func (h *Handler) Get(obj interface{}) (*unstructured.Unstructured, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.createUnstructured(val.(*unstructured.Unstructured))
		}
		return h.GetFromObject(val)
	case *unstructured.Unstructured:
		return h.getUnstructured(val)
	case unstructured.Unstructured:
		return h.getUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	default:
		return nil, ErrInvalidType
	}
}

// GetByName gets unstructured k8s resource with given name.
func (h *Handler) GetByName(name string) (*unstructured.Unstructured, error) {
	if h.IsNamespacedResource() {
		return h.dynamicClient.Resource(h.gvr).Namespace(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
	}
	return h.dynamicClient.Resource(h.gvr).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets unstructured k8s resource from yaml file.
func (h *Handler) GetFromFile(filename string) (*unstructured.Unstructured, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets unstructured k8s resource from bytes.
func (h *Handler) GetFromBytes(data []byte) (*unstructured.Unstructured, error) {
	unstructJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	unstructObj := &unstructured.Unstructured{}
	if err = json.Unmarshal(unstructJson, unstructObj); err != nil {
		return nil, err
	}
	return h.getUnstructured(unstructObj)
}

// GetFromObject gets unstructured k8s resource from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*unstructured.Unstructured, error) {
	unstructMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	return h.getUnstructured(&unstructured.Unstructured{Object: unstructMap})
}

// GetFromMap gets unstructured k8s resource from map[string]interface{}.
func (h *Handler) GetFromMap(obj map[string]interface{}) (*unstructured.Unstructured, error) {
	return h.getUnstructured(&unstructured.Unstructured{Object: obj})
}

// getUnstructured
func (h *Handler) getUnstructured(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	if h.IsNamespacedResource() {
		return h.dynamicClient.Resource(h.gvr).Namespace(h.namespace).Get(h.ctx, obj.GetName(), h.Options.GetOptions)
	}
	return h.dynamicClient.Resource(h.gvr).Get(h.ctx, obj.GetName(), h.Options.GetOptions)
}
