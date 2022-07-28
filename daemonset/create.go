package daemonset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates daemonset from type string, []byte, *appsv1.DaemonSet,
// appsv1.DaemonSet, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*appsv1.DaemonSet, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *appsv1.DaemonSet:
		return h.CreateFromObject(val)
	case appsv1.DaemonSet:
		return h.CreateFromObject(&val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates daemonset from yaml file.
func (h *Handler) CreateFromFile(filename string) (*appsv1.DaemonSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates daemonset from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*appsv1.DaemonSet, error) {
	dsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ds := &appsv1.DaemonSet{}
	err = json.Unmarshal(dsJson, ds)
	if err != nil {
		return nil, err
	}
	return h.createDaemonset(ds)
}

// CreateFromObject creates daemonset from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*appsv1.DaemonSet, error) {
	ds, ok := obj.(*appsv1.DaemonSet)
	if !ok {
		return nil, fmt.Errorf("object is not *appsv1.DaemonSet")
	}
	return h.createDaemonset(ds)
}

// CreateFromUnstructured creates daemonset from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*appsv1.DaemonSet, error) {
	ds := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ds)
	if err != nil {
		return nil, err
	}
	return h.createDaemonset(ds)
}

// CreateFromMap creates daemonset from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*appsv1.DaemonSet, error) {
	ds := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ds)
	if err != nil {
		return nil, err
	}
	return h.createDaemonset(ds)
}

// createDaemonset
func (h *Handler) createDaemonset(ds *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	var namespace string
	if len(ds.Namespace) != 0 {
		namespace = ds.Namespace
	} else {
		namespace = h.namespace
	}
	ds.ResourceVersion = ""
	ds.UID = ""
	return h.clientset.AppsV1().DaemonSets(namespace).Create(h.ctx, ds, h.Options.CreateOptions)
}
