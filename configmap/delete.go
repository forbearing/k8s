package configmap

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete configmap from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	cmJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	configmap := &corev1.ConfigMap{}
	err = json.Unmarshal(cmJson, configmap)
	if err != nil {
		return err
	}

	var namespace string
	if len(configmap.Namespace) != 0 {
		namespace = configmap.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(configmap.Name)
}

// DeleteFromFile delete configmap from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete configmap by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().ConfigMaps(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete configmap by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
