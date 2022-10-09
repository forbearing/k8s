package configmap

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

// Create creates configmap from type string, []byte, *corev1.ConfigMap,
// corev1.ConfigMap, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*corev1.ConfigMap, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *corev1.ConfigMap:
		return h.CreateFromObject(val)
	case corev1.ConfigMap:
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

// CreateFromFile creates configmap from yaml or json file.
func (h *Handler) CreateFromFile(filename string) (*corev1.ConfigMap, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates configmap from bytes data.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.ConfigMap, error) {
	cmJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cm := &corev1.ConfigMap{}
	if err = json.Unmarshal(cmJson, cm); err != nil {
		return nil, err
	}
	return h.createConfigmap(cm)
}

// CreateFromObject creates configmap from metav1.Object or runtime.Object.
func (h *Handler) CreateFromObject(obj interface{}) (*corev1.ConfigMap, error) {
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.ConfigMap")
	}
	return h.createConfigmap(cm)
}

// CreateFromUnstructured creates configmap from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cm)
	if err != nil {
		return nil, err
	}
	return h.createConfigmap(cm)
}

// CreateFromMap creates configmap from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cm)
	if err != nil {
		return nil, err
	}
	return h.createConfigmap(cm)
}

// createConfigmap
func (h *Handler) createConfigmap(cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	namespace := cm.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	cm.ResourceVersion = ""
	cm.UID = ""
	return h.clientset.CoreV1().ConfigMaps(namespace).Create(h.ctx, cm, h.Options.CreateOptions)
}
