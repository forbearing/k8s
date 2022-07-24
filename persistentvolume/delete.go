package persistentvolume

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Delete deletes persistentvolume from type string, []byte, *corev1.PersistentVolume,
// corev1.PersistentVolume, runtime.Object or map[string]interface{}.

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
	case runtime.Object:
		return h.DeleteFromObject(val)
	case map[string]interface{}:
		return h.DeleteFromUnstructured(val)
	default:
		return ERR_TYPE_DELETE
	}
}

// DeleteByName deletes persistentvolume by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().PersistentVolumes().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes persistentvolume from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes persistentvolume from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	pvJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	pv := &corev1.PersistentVolume{}
	err = json.Unmarshal(pvJson, pv)
	if err != nil {
		return err
	}
	return h.deletePV(pv)
}

// DeleteFromObject deletes persistentvolume from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	pv, ok := obj.(*corev1.PersistentVolume)
	if !ok {
		return fmt.Errorf("object is not *corev1.PersistentVolume")
	}
	return h.deletePV(pv)
}

// DeleteFromUnstructured deletes persistentvolume from map[string]interface{}.
func (h *Handler) DeleteFromUnstructured(u map[string]interface{}) error {
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
