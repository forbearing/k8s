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

// Create creates persistentvolume from type string, []byte, *corev1.PersistentVolume,
// corev1.PersistentVolume, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*corev1.PersistentVolume, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *corev1.PersistentVolume:
		return h.CreateFromObject(val)
	case corev1.PersistentVolume:
		return h.CreateFromObject(&val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.CreateFromObject(val)
	default:
		return nil, ErrInvalidCreateType
	}
}

// CreateFromFile creates persistentvolume from yaml or json file.
func (h *Handler) CreateFromFile(filename string) (*corev1.PersistentVolume, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates persistentvolume from bytes data.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.PersistentVolume, error) {
	pvJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pv := &corev1.PersistentVolume{}
	if err = json.Unmarshal(pvJson, pv); err != nil {
		return nil, err
	}
	return h.createPV(pv)
}

// CreateFromObject creates persistentvolume from metav1.Object or runtime.Object.
func (h *Handler) CreateFromObject(obj interface{}) (*corev1.PersistentVolume, error) {
	pv, ok := obj.(*corev1.PersistentVolume)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.PersistentVolume")
	}
	return h.createPV(pv)
}

// CreateFromUnstructured creates persistentvolume from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*corev1.PersistentVolume, error) {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pv)
	if err != nil {
		return nil, err
	}
	return h.createPV(pv)
}

// CreateFromMap creates persistentvolume from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*corev1.PersistentVolume, error) {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pv)
	if err != nil {
		return nil, err
	}
	return h.createPV(pv)
}

// createPV
func (h *Handler) createPV(pv *corev1.PersistentVolume) (*corev1.PersistentVolume, error) {
	pv.ResourceVersion = ""
	pv.UID = ""
	return h.clientset.CoreV1().PersistentVolumes().Create(h.ctx, pv, h.Options.CreateOptions)
}
