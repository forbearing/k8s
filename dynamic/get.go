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

// Get gets unstructured k8s resource from type string, []byte, metav1.Object,
// runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.
//
// If psssed parameter type is string, it will call GetByName insteard of GetFromFile.
// You should always explicitly call GetFromFile to delete a unstructured object
// from filename.
//
// GetByName requires WithGVK() to explicitly specify the k8s resource's GroupVersionKind.
// GetFromFile, GetFromBytes and GetFromMap will find GVK and GVR from
// the provided structured or unstructured data, it's not reuqired to call WithGVK().
func (h *Handler) Get(obj interface{}) (*unstructured.Unstructured, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *unstructured.Unstructured:
		return h.getUnstructured(val)
	case unstructured.Unstructured:
		return h.getUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	case metav1.Object, runtime.Object:
		return h.GetFromObject(val)
	default:
		return nil, ErrInvalidGetType
	}
}

// GetByName gets unstructured k8s resource with given name.
func (h *Handler) GetByName(name string) (*unstructured.Unstructured, error) {
	var err error
	if h.gvr, err = utilrestmapper.GVKToGVR(h.restMapper, h.gvk); err != nil {
		return nil, err
	}
	if h.isNamespaced, err = utilrestmapper.IsNamespaced(h.restMapper, h.gvk); err != nil {
		return nil, err
	}
	if h.gvk.Kind == types.KindJob || h.gvk.Kind == types.KindCronJob {
		h.SetPropagationPolicy("background")
	}

	if h.isNamespaced {
		return h.dynamicClient.Resource(h.gvr).Namespace(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
	}
	return h.dynamicClient.Resource(h.gvr).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets unstructured k8s resource from yaml or json file.
func (h *Handler) GetFromFile(filename string) (*unstructured.Unstructured, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets unstructured k8s resource from bytes data.
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

// GetFromObject gets unstructured k8s resource from metav1.Object or runtime.Object.
func (h *Handler) GetFromObject(obj interface{}) (*unstructured.Unstructured, error) {
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

	if h.isNamespaced {
		namespace := obj.GetNamespace()
		if len(namespace) == 0 {
			namespace = h.namespace
		}
		return h.dynamicClient.Resource(h.gvr).Namespace(namespace).Get(h.ctx, obj.GetName(), h.Options.GetOptions)
	}
	return h.dynamicClient.Resource(h.gvr).Get(h.ctx, obj.GetName(), h.Options.GetOptions)
}
