package configmap

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update configmap from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*corev1.ConfigMap, error) {
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

	return h.clientset.CoreV1().ConfigMaps(namespace).Update(h.ctx, configmap, h.Options.UpdateOptions)
}

// UpdateFromBytes update configmap from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.ConfigMap, error) {
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

	return h.clientset.CoreV1().ConfigMaps(namespace).Update(h.ctx, configmap, h.Options.UpdateOptions)
}

// UpdateFromFile update configmap from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.ConfigMap, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update configmap from file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*corev1.ConfigMap, error) {
	return h.UpdateFromFile(filename)
}
