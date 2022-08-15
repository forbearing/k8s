package persistentvolume

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets persistentvolume from type string, []byte, *corev1.PersistentVolume,
// corev1.PersistentVolume, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a persistentvolume from file path.
func (h *Handler) Get(obj interface{}) (*corev1.PersistentVolume, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *corev1.PersistentVolume:
		return h.GetFromObject(val)
	case corev1.PersistentVolume:
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
		return nil, ErrInvalidGetType
	}
}

// GetByName gets persistentvolume by name.
func (h *Handler) GetByName(name string) (*corev1.PersistentVolume, error) {
	return h.clientset.CoreV1().PersistentVolumes().Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets persistentvolume from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.PersistentVolume, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets persistentvolume from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.PersistentVolume, error) {
	pvJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pv := &corev1.PersistentVolume{}
	if err = json.Unmarshal(pvJson, pv); err != nil {
		return nil, err
	}
	return h.getPV(pv)
}

// GetFromObject gets persistentvolume from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*corev1.PersistentVolume, error) {
	pv, ok := obj.(*corev1.PersistentVolume)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.PersistentVolume")
	}
	return h.getPV(pv)
}

// GetFromUnstructured gets persistentvolume from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*corev1.PersistentVolume, error) {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pv)
	if err != nil {
		return nil, err
	}
	return h.getPV(pv)
}

// GetFromMap gets persistentvolume from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*corev1.PersistentVolume, error) {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pv)
	if err != nil {
		return nil, err
	}
	return h.getPV(pv)
}

// getPV
// It's necessary to get a new persistentvolume resource from a old persistentvolume resource,
// because old persistentvolume usually don't have persistentvolume.Status field.
func (h *Handler) getPV(pv *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	return h.clientset.CoreV1().PersistentVolumes().Get(h.ctx, pv.Name, h.Options.GetOptions)
}
