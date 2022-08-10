package dynamic

import (
	"encoding/json"
	"io/ioutil"
	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates unstructured k8s resource from type string, []byte,
// runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*unstructured.Unstructured, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.updateUnstructured(val.(*unstructured.Unstructured))
		}
		return h.UpdateFromObject(val)
	case *unstructured.Unstructured:
		return h.updateUnstructured(val)
	case unstructured.Unstructured:
		return h.updateUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	default:
		return nil, ErrInvalidType
	}
}

// UpdateFromFile updates unstructured k8s resource from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*unstructured.Unstructured, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates unstructured k8s resource from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*unstructured.Unstructured, error) {
	unstructJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	unstructObj := &unstructured.Unstructured{}
	if err = json.Unmarshal(unstructJson, unstructObj); err != nil {
		return nil, err
	}
	return h.updateUnstructured(unstructObj)
}

// UpdateFromObject updates unstructured k8s resource from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*unstructured.Unstructured, error) {
	unstructMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	return h.updateUnstructured(&unstructured.Unstructured{Object: unstructMap})
}

// UpdateFromMap updates unstructured k8s resource from map[string]interface{}.
func (h *Handler) UpdateFromMap(obj map[string]interface{}) (*unstructured.Unstructured, error) {
	return h.updateUnstructured(&unstructured.Unstructured{Object: obj})
}

// updateUnstructured
func (h *Handler) updateUnstructured(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	obj.SetUID("")
	obj.SetResourceVersion("")
	if h.IsNamespacedResource() {
		return h.dynamicClient.Resource(h.gvr).Namespace(h.namespace).Update(h.ctx, obj, h.Options.UpdateOptions)
	}
	return h.dynamicClient.Resource(h.gvr).Update(h.ctx, obj, h.Options.UpdateOptions)
}
