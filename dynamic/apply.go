package dynamic

import (
	"encoding/json"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Apply applies unstructured k8s resource from type string, []byte,
// runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.
//
// It's not necessary to explicitly specify the GVK or GVR by calling WithGVK(),
// Apply() will find the GVK and GVR by RESTMapper and apply the k8s resource
// that defined in yaml file, json file, bytes data, map[string]interface{}
// or runtime.Object.
func (h *Handler) Apply(obj interface{}) (*unstructured.Unstructured, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *unstructured.Unstructured:
		return h.applyUnstructured(val)
	case unstructured.Unstructured:
		return h.applyUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	case runtime.Object:
		return h.ApplyFromObject(val)
	default:
		return nil, ErrInvalidApplyType
	}
}

// ApplyFromFile applies unstructured k8s resource from yaml file.
func (h *Handler) ApplyFromFile(filename string) (*unstructured.Unstructured, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.ApplyFromBytes(data)
}

// ApplyFromBytes applies unstructured k8s resource from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (*unstructured.Unstructured, error) {
	unstructJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	unstructObj := &unstructured.Unstructured{}
	if err = json.Unmarshal(unstructJson, unstructObj); err != nil {
		return nil, err
	}
	return h.applyUnstructured(unstructObj)
}

// ApplyFromObject applies unstructured k8s resource from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*unstructured.Unstructured, error) {
	unstructMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	return h.applyUnstructured(&unstructured.Unstructured{Object: unstructMap})
}

// ApplyFromMap applies unstructured k8s resource from map[string]interface{}.
func (h *Handler) ApplyFromMap(obj map[string]interface{}) (*unstructured.Unstructured, error) {
	return h.applyUnstructured(&unstructured.Unstructured{Object: obj})
}

// applyUnstructured
func (h *Handler) applyUnstructured(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	_, err := h.createUnstructured(obj)
	if errors.IsAlreadyExists(err) {
		return h.Update(obj)
	}
	return obj, err
}
