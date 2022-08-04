package daemonset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets daemonset from type string, []byte, *appsv1.DaemonSet,
// appsv1.DaemonSet, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a daemonset from file path.
func (h *Handler) Get(obj interface{}) (*appsv1.DaemonSet, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *appsv1.DaemonSet:
		return h.GetFromObject(val)
	case appsv1.DaemonSet:
		return h.GetFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.GetFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.GetFromObject(val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets daemonset by name.
func (h *Handler) GetByName(name string) (*appsv1.DaemonSet, error) {
	return h.clientset.AppsV1().DaemonSets(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets daemonset from yaml file.
func (h *Handler) GetFromFile(filename string) (*appsv1.DaemonSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets daemonset from bytes.
func (h *Handler) GetFromBytes(data []byte) (*appsv1.DaemonSet, error) {
	dsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ds := &appsv1.DaemonSet{}
	if err = json.Unmarshal(dsJson, ds); err != nil {
		return nil, err
	}
	return h.getDaemonset(ds)
}

// GetFromObject gets daemonset from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*appsv1.DaemonSet, error) {
	ds, ok := obj.(*appsv1.DaemonSet)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.DaemonSet")
	}
	return h.getDaemonset(ds)
}

// GetFromUnstructured gets daemonset from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*appsv1.DaemonSet, error) {
	ds := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ds)
	if err != nil {
		return nil, err
	}
	return h.getDaemonset(ds)
}

// GetFromMap gets daemonset from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*appsv1.DaemonSet, error) {
	ds := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ds)
	if err != nil {
		return nil, err
	}
	return h.getDaemonset(ds)
}

// getDaemonset
// It's necessary to get a new daemonset resource from a old daemonset resource,
// because old daemonset usually don't have daemonset.Status field.
func (h *Handler) getDaemonset(ds *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	var namespace string
	if len(ds.Namespace) != 0 {
		namespace = ds.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.AppsV1().DaemonSets(namespace).Get(h.ctx, ds.Name, h.Options.GetOptions)
}
