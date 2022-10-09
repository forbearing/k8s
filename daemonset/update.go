package daemonset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates daemonset from type string, []byte, *appsv1.DaemonSet,
// appsv1.DaemonSet, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*appsv1.DaemonSet, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *appsv1.DaemonSet:
		return h.UpdateFromObject(val)
	case appsv1.DaemonSet:
		return h.UpdateFromObject(&val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.UpdateFromObject(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates daemonset from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*appsv1.DaemonSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates daemonset from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*appsv1.DaemonSet, error) {
	dsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ds := &appsv1.DaemonSet{}
	if err = json.Unmarshal(dsJson, ds); err != nil {
		return nil, err
	}
	return h.updateDaemonset(ds)
}

// UpdateFromObject updates daemonset from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*appsv1.DaemonSet, error) {
	ds, ok := obj.(*appsv1.DaemonSet)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.DaemonSet")
	}
	return h.updateDaemonset(ds)
}

// UpdateFromUnstructured updates daemonset from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*appsv1.DaemonSet, error) {
	ds := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ds)
	if err != nil {
		return nil, err
	}
	return h.updateDaemonset(ds)
}

// UpdateFromMap updates daemonset from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*appsv1.DaemonSet, error) {
	ds := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ds)
	if err != nil {
		return nil, err
	}
	return h.updateDaemonset(ds)
}

// updateDaemonset
func (h *Handler) updateDaemonset(ds *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	namespace := ds.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	ds.ResourceVersion = ""
	ds.UID = ""
	return h.clientset.AppsV1().DaemonSets(namespace).Update(h.ctx, ds, h.Options.UpdateOptions)
}
