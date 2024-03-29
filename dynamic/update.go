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

// Update updates unstructured k8s resource from type string, []byte, metav1.Object,
// runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.
//
// It's not necessary to explicitly specify the GVK or GVR by calling WithGVK(),
// Update() will find the GVK and GVR by RESTMapper and update the k8s resource
// that defined in yaml file, json file, bytes data, map[string]interface{}
// or runtime.Object.
func (h *Handler) Update(obj interface{}) (*unstructured.Unstructured, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *unstructured.Unstructured:
		return h.updateUnstructured(val)
	case unstructured.Unstructured:
		return h.updateUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.UpdateFromObject(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates unstructured k8s resource from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*unstructured.Unstructured, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates unstructured k8s resource from bytes data.
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

// UpdateFromObject updates unstructured k8s resource from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*unstructured.Unstructured, error) {
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
		return h.dynamicClient.Resource(h.gvr).Namespace(namespace).Update(h.ctx, obj, h.Options.UpdateOptions)
	}
	return h.dynamicClient.Resource(h.gvr).Update(h.ctx, obj, h.Options.UpdateOptions)
}
