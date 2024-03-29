package persistentvolume

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Delete deletes persistentvolume from type string, []byte, *corev1.PersistentVolume,
// corev1.PersistentVolume, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a persistentvolume from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *corev1.PersistentVolume:
		return h.DeleteFromObject(val)
	case corev1.PersistentVolume:
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

// DeleteByName deletes persistentvolume by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().PersistentVolumes().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes persistentvolume from yaml or json file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes persistentvolume from bytes data.
func (h *Handler) DeleteFromBytes(data []byte) error {
	pvJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	pv := &corev1.PersistentVolume{}
	if err = json.Unmarshal(pvJson, pv); err != nil {
		return err
	}
	return h.deletePV(pv)
}

// DeleteFromObject deletes persistentvolume from metav1.Object or runtime.Object.
func (h *Handler) DeleteFromObject(obj interface{}) error {
	pv, ok := obj.(*corev1.PersistentVolume)
	if !ok {
		return fmt.Errorf("object type is not *corev1.PersistentVolume")
	}
	return h.deletePV(pv)
}

// DeleteFromUnstructured deletes persistentvolume from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pv)
	if err != nil {
		return err
	}
	return h.deletePV(pv)
}

// DeleteFromMap deletes persistentvolume from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pv)
	if err != nil {
		return err
	}
	return h.deletePV(pv)
}

// deletePV
func (h *Handler) deletePV(pv *corev1.PersistentVolume) error {
	return h.clientset.CoreV1().PersistentVolumes().Delete(h.ctx, pv.Name, h.Options.DeleteOptions)
}
