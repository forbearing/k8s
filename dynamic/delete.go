package dynamic

import (
	"encoding/json"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Delete deletes unstructured k8s resource from type string, []byte,
// runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.

// If psssed parameter type is string, it will call DeleteByName insteard of DeleteFromFile.
// You  should always explicitly call DeleteFromFile to delete a unstructured object
// from filename.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case runtime.Object:
		return h.DeleteFromObject(val)
	case *unstructured.Unstructured:
		return h.deleteUnstructured(val)
	case unstructured.Unstructured:
		return h.deleteUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	default:
		return ErrInvalidType
	}
}

// DeleteByName deletes unstructured k8s resource with given name.
func (h *Handler) DeleteByName(name string) error {
	if h.IsNamespacedResource() {
		//logrus.Info("namespacedResource")
		return h.dynamicClient.Resource(h.gvr()).Namespace(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
	}
	//logrus.Info("not namespacedResource")
	return h.dynamicClient.Resource(h.gvr()).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes unstructured k8s resource from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes unstructured k8s resource from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	unstructJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	unstructObj := &unstructured.Unstructured{}
	if err = json.Unmarshal(unstructJson, unstructObj); err != nil {
		return err
	}
	return h.deleteUnstructured(unstructObj)
}

// DeleteFromObject deletes unstructured k8s resource from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	unstructMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return err
	}
	return h.deleteUnstructured(&unstructured.Unstructured{Object: unstructMap})
}

// DeleteFromMap deletes unstructured k8s resource from map[string]interface{}.
func (h *Handler) DeleteFromMap(obj map[string]interface{}) error {
	return h.deleteUnstructured(&unstructured.Unstructured{Object: obj})
}

// deleteUnstructured
func (h *Handler) deleteUnstructured(obj *unstructured.Unstructured) error {
	if h.IsNamespacedResource() {
		//logrus.Info("namespacedResource")
		return h.dynamicClient.Resource(h.gvr()).Namespace(h.namespace).Delete(h.ctx, obj.GetName(), h.Options.DeleteOptions)
	}
	//logrus.Info("not namespacedResource")
	return h.dynamicClient.Resource(h.gvr()).Delete(h.ctx, obj.GetName(), h.Options.DeleteOptions)
}
