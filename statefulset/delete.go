package statefulset

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

// Delete deletes statefulset from type string, []byte, *appsv1.StatefulSet,
// appsv1.StatefulSet, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a statefulset from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *appsv1.StatefulSet:
		return h.DeleteFromObject(val)
	case appsv1.StatefulSet:
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

// DeleteByName deletes statefulset by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.AppsV1().StatefulSets(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes statefulset from yaml or json file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes statefulset from bytes data.
func (h *Handler) DeleteFromBytes(data []byte) error {
	stsJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	sts := &appsv1.StatefulSet{}
	if err = json.Unmarshal(stsJson, sts); err != nil {
		return err
	}
	return h.deleteStatefulset(sts)
}

// DeleteFromObject deletes statefulset from metav1.Object or runtime.Object.
func (h *Handler) DeleteFromObject(obj interface{}) error {
	sts, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		return fmt.Errorf("object type is not *appsv1.StatefulSet")
	}
	return h.deleteStatefulset(sts)
}

// DeleteFromUnstructured deletes statefulset from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sts)
	if err != nil {
		return err
	}
	return h.deleteStatefulset(sts)
}

// DeleteFromMap deletes statefulset from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sts)
	if err != nil {
		return err
	}
	return h.deleteStatefulset(sts)
}

// deleteStatefulset
func (h *Handler) deleteStatefulset(sts *appsv1.StatefulSet) error {
	namespace := sts.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.AppsV1().StatefulSets(namespace).Delete(h.ctx, sts.Name, h.Options.DeleteOptions)
}
