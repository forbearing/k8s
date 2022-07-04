package configmap

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get configmap from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.ConfigMap, error) {
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

	return h.WithNamespace(namespace).GetByName(configmap.Name)
}

// GetFromFile get configmap from yaml file.
func (h *Handler) GetFromFile(filename string) (configmap *corev1.ConfigMap, err error) {
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	configmap, err = h.GetFromBytes(data)
	return
}

// GetByName get configmap by name.
func (h *Handler) GetByName(name string) (*corev1.ConfigMap, error) {
	return h.clientset.CoreV1().ConfigMaps(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get configmap by name.
func (h *Handler) Get(name string) (*corev1.ConfigMap, error) {
	return h.GetByName(name)
}
