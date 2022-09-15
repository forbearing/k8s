package dynamic

import (
	"encoding/json"
	"io/ioutil"

	"github.com/forbearing/k8s/types"
	utilrestmapper "github.com/forbearing/k8s/util/restmapper"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates unstructured k8s resource from type string, []byte,
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
	case runtime.Object:
		return h.CreateFromObject(val)
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
	var (
		err          error
		gvk          schema.GroupVersionKind
		gvr          schema.GroupVersionResource
		isNamespaced bool
	)
	if gvr, err = utilrestmapper.FindGVR(h.restMapper, obj); err != nil {
		logrus.Error("get gvr failed", err)
		return nil, err
	}
	if gvk, err = utilrestmapper.FindGVK(h.restMapper, obj); err != nil {
		logrus.Error("get gvk failed", err)
		return nil, err
	}
	if isNamespaced, err = utilrestmapper.IsNamespaced(h.restMapper, gvk); err != nil {
		logrus.Error("get isNamespaced failed", err)
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
		return h.dynamicClient.Resource(gvr).Namespace(namespace).Create(h.ctx, obj, h.Options.CreateOptions)
	}
	return h.dynamicClient.Resource(gvr).Create(h.ctx, obj, h.Options.CreateOptions)
}
