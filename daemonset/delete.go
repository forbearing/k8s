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

// Delete deletes daemonset from type string, []byte, *appsv1.DaemonSet,
// appsv1.DaemonSet, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a daemonset from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *appsv1.DaemonSet:
		return h.DeleteFromObject(val)
	case appsv1.DaemonSet:
		return h.DeleteFromObject(&val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	case metav1.Object, runtime.Object:
		return h.DeleteFromObject(val)
	default:
		return ErrInvalidDeleteType
	}
}

// DeleteByName deletes daemonset by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.AppsV1().DaemonSets(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes daemonset from yaml or json file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes daemonset from bytes data.
func (h *Handler) DeleteFromBytes(data []byte) error {
	dsJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	ds := &appsv1.DaemonSet{}
	if err = json.Unmarshal(dsJson, ds); err != nil {
		return err
	}
	return h.deleteDaemonset(ds)
}

// DeleteFromObject deletes daemonset from metav1.Object or runtime.Object.
func (h *Handler) DeleteFromObject(obj interface{}) error {
	ds, ok := obj.(*appsv1.DaemonSet)
	if !ok {
		return fmt.Errorf("object type is not *appsv1.DaemonSet")
	}
	return h.deleteDaemonset(ds)
}

// DeleteFromUnstructured deletes daemonset from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	ds := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ds)
	if err != nil {
		return err
	}
	return h.deleteDaemonset(ds)
}

// DeleteFromMap deletes daemonset from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	ds := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ds)
	if err != nil {
		return err
	}
	return h.deleteDaemonset(ds)
}

// deleteDaemonset
func (h *Handler) deleteDaemonset(ds *appsv1.DaemonSet) error {
	namespace := ds.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.AppsV1().DaemonSets(namespace).Delete(h.ctx, ds.Name, h.Options.DeleteOptions)
}
