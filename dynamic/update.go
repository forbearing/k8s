package dynamic

import (
	"encoding/json"
	"io/ioutil"

	"github.com/forbearing/k8s/types"
	utilrestmapper "github.com/forbearing/k8s/util/restmapper"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates unstructured k8s resource from type string, []byte,
// runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.
//
// It's not necessary to explicitly specify the GVK or GVR, Update() will find
// the GVK and GVR by RESTMapper and update the k8s resource that defined in
// yaml file, json file, bytes data, map[string]interface{} or runtime.Object.
func (h *Handler) Update(obj interface{}) (*unstructured.Unstructured, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	//case runtime.Object:
	//    if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
	//        return h.updateUnstructured(val.(*unstructured.Unstructured))
	//    }
	//    return h.UpdateFromObject(val)
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
	var (
		err          error
		gvk          schema.GroupVersionKind
		gvr          schema.GroupVersionResource
		isNamespaced bool
	)
	if gvr, err = utilrestmapper.FindGVR(h.restMapper, obj); err != nil {
		return nil, err
	}
	if gvk, err = utilrestmapper.FindGVK(h.restMapper, obj); err != nil {
		return nil, err
	}
	if isNamespaced, err = utilrestmapper.IsNamespaced(h.restMapper, gvk); err != nil {
		return nil, err
	}
	if gvk.Kind == types.KindJob || gvk.Kind == types.KindCronJob {
		h.SetPropagationPolicy("background")
	}

	obj.SetUID("")
	obj.SetResourceVersion("")
	if isNamespaced {
		var namespace string
		if len(obj.GetNamespace()) != 0 {
			namespace = obj.GetNamespace()
		} else {
			namespace = h.namespace
		}
		return h.dynamicClient.Resource(gvr).Namespace(namespace).Update(h.ctx, obj, h.Options.UpdateOptions)
	}
	return h.dynamicClient.Resource(gvr).Update(h.ctx, obj, h.Options.UpdateOptions)
}
