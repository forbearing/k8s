package configmap

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create configmap from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*corev1.ConfigMap, error) {
	configmap := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, configmap)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(configmap.Namespace) != 0 {
		namespace = configmap.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.CoreV1().ConfigMaps(namespace).Create(h.ctx, configmap, h.Options.CreateOptions)
}

// CreateFromBytes create configmap from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.ConfigMap, error) {
	cmJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	configmap := &corev1.ConfigMap{}
	err = json.Unmarshal(cmJson, configmap)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(configmap.Namespace) != 0 {
		namespace = configmap.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.CoreV1().ConfigMaps(namespace).Create(h.ctx, configmap, h.Options.CreateOptions)
}

// CreateFromFile create configmap from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.ConfigMap, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create configmap from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*corev1.ConfigMap, error) {
	return h.CreateFromFile(filename)
}
