package dynamic

import (
	"encoding/json"
	"io/ioutil"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates unstructured k8s resource from type string, []byte,
// runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*unstructured.Unstructured, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case *unstructured.Unstructured:
		return h.createUnstructured(val)
	case unstructured.Unstructured:
		return h.createUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	default:
		return nil, ErrInvalidType
	}
}

// CreateFromFile creates unstructured k8s resource from yaml file.
func (h *Handler) CreateFromFile(filename string) (*unstructured.Unstructured, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates unstructured k8s resource from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*unstructured.Unstructured, error) {
	unstructJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	unstructObj := &unstructured.Unstructured{}
	if err = json.Unmarshal(unstructJson, unstructObj); err != nil {
		return nil, err
	}
	return h.createUnstructured(unstructObj)
}

// CreateFromObject creates unstructured k8s resource from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*unstructured.Unstructured, error) {
	unstructMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	return h.createUnstructured(&unstructured.Unstructured{Object: unstructMap})
}

// CreateFromMap creates unstructured k8s resource from map[string]interface{}.
func (h *Handler) CreateFromMap(obj map[string]interface{}) (*unstructured.Unstructured, error) {
	return h.createUnstructured(&unstructured.Unstructured{Object: obj})
}

// createUnstructured
func (h *Handler) createUnstructured(obj *unstructured.Unstructured) (*unstructured.Unstructured, error) {
	obj.SetUID("")
	obj.SetResourceVersion("")
	//logrus.Info(h.gvr())
	if h.IsNamespacedResource() {
		//logrus.Info("namespacedResource")
		return h.dynamicClient.Resource(h.gvr()).Namespace(h.namespace).Create(h.ctx, obj, h.Options.CreateOptions)
	}
	//logrus.Info("not namespacedResource")
	return h.dynamicClient.Resource(h.gvr()).Create(h.ctx, obj, h.Options.CreateOptions)
}
