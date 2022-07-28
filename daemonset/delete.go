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

// Delete deletes daemonset from type string, []byte, *appsv1.DaemonSet,
// appsv1.DaemonSet, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

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
	case runtime.Object:
		return h.DeleteFromObject(val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	default:
		return ERR_TYPE_DELETE
	}
}

// DeleteByName deletes daemonset by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.AppsV1().DaemonSets(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes daemonset from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes daemonset from bytes.
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

// DeleteFromObject deletes daemonset from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
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
	var namespace string
	if len(ds.Namespace) != 0 {
		namespace = ds.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.AppsV1().DaemonSets(namespace).Delete(h.ctx, ds.Name, h.Options.DeleteOptions)
}
