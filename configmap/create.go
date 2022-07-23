package configmap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates configmap from type string, []byte, *corev1.ConfigMap,
// corev1.ConfigMap, runtime.Object or map[string]interface{}.
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
	case runtime.Object:
		return h.CreateFromObject(val)
	case map[string]interface{}:
		return h.CreateFromUnstructured(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates configmap from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.ConfigMap, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates configmap from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.ConfigMap, error) {
	cmJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cm := &corev1.ConfigMap{}
	err = json.Unmarshal(cmJson, cm)
	if err != nil {
		return nil, err
	}
	return h.createConfigmap(cm)
}

// CreateFromObject creates configmap from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*corev1.ConfigMap, error) {
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("object is not *corev1.ConfigMap")
	}
	return h.createConfigmap(cm)
}

// CreateFromUnstructured creates configmap from map[string]interface{}.
func (h *Handler) CreateFromUnstructured(u map[string]interface{}) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cm)
	if err != nil {
		return nil, err
	}
	return h.createConfigmap(cm)
}

// createConfigmap
func (h *Handler) createConfigmap(cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	var namespace string
	if len(cm.Namespace) != 0 {
		namespace = cm.Namespace
	} else {
		namespace = h.namespace
	}
	cm.ResourceVersion = ""
	cm.UID = ""
	return h.clientset.CoreV1().ConfigMaps(namespace).Create(h.ctx, cm, h.Options.CreateOptions)
}
