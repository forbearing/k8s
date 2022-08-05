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

// Update updates persistentvolume from type string, []byte, *corev1.PersistentVolume,
// corev1.PersistentVolume, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*corev1.PersistentVolume, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *corev1.PersistentVolume:
		return h.UpdateFromObject(val)
	case corev1.PersistentVolume:
		return h.UpdateFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.UpdateFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.UpdateFromObject(val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates persistentvolume from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.PersistentVolume, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates persistentvolume from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.PersistentVolume, error) {
	pvJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pv := &corev1.PersistentVolume{}
	if err = json.Unmarshal(pvJson, pv); err != nil {
		return nil, err
	}
	return h.updatePV(pv)
}

// UpdateFromObject updates persistentvolume from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*corev1.PersistentVolume, error) {
	pv, ok := obj.(*corev1.PersistentVolume)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.PersistentVolume")
	}
	return h.updatePV(pv)
}

// UpdateFromUnstructured updates persistentvolume from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*corev1.PersistentVolume, error) {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pv)
	if err != nil {
		return nil, err
	}
	return h.updatePV(pv)
}

// UpdateFromMap updates persistentvolume from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*corev1.PersistentVolume, error) {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pv)
	if err != nil {
		return nil, err
	}
	return h.updatePV(pv)
}

// updatePV
func (h *Handler) updatePV(pv *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	pv.ResourceVersion = ""
	pv.UID = ""
	return h.clientset.CoreV1().PersistentVolumes().Update(h.ctx, pv, h.Options.UpdateOptions)
}
