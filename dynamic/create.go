package dynamic

import (
	"encoding/json"
	"io/ioutil"

	"github.com/forbearing/k8s/types"
	utilrestmapper "github.com/forbearing/k8s/util/restmapper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates unstructured k8s resource from type string, []byte, metav1.Object,
// runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.
//
// It's not necessary to explicitly specify the GVK or GVR by calling WithGVK(),
// Create() will find the GVK and GVR by RESTMapper and create the k8s resource
// that defined in yaml file, json file, bytes data, map[string]interface{}
// or runtime.Object.
func (h *Handler) Create(obj interface{}) (*unstructured.Unstructured, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *unstructured.Unstructured:
		return h.createUnstructured(val)
	case unstructured.Unstructured:
		return h.createUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.CreateFromObject(val)
	default:
		return nil, ErrInvalidCreateType
	}
}

// CreateFromFile creates unstructured k8s resource from yaml or json file.
func (h *Handler) CreateFromFile(filename string) (*unstructured.Unstructured, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates unstructured k8s resource from bytes data.
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

// CreateFromObject creates unstructured k8s resource from metav1.Object or runtime.Object.
func (h *Handler) CreateFromObject(obj interface{}) (*unstructured.Unstructured, error) {
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
	var err error
	if h.gvr, err = utilrestmapper.FindGVR(h.restMapper, obj); err != nil {
		return nil, err
	}
	if h.gvk, err = utilrestmapper.FindGVK(h.restMapper, obj); err != nil {
		return nil, err
	}
	if h.isNamespaced, err = utilrestmapper.IsNamespaced(h.restMapper, h.gvk); err != nil {
		return nil, err
	}
	if h.gvk.Kind == types.KindJob || h.gvk.Kind == types.KindCronJob {
		h.SetPropagationPolicy("background")
	}

	obj.SetUID("")
	obj.SetResourceVersion("")
	if h.isNamespaced {
		namespace := obj.GetNamespace()
		if len(namespace) == 0 {
			namespace = h.namespace
		}
		return h.dynamicClient.Resource(h.gvr).Namespace(namespace).Create(h.ctx, obj, h.Options.CreateOptions)
	}
	return h.dynamicClient.Resource(h.gvr).Create(h.ctx, obj, h.Options.CreateOptions)
}
