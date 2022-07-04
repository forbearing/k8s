package namespace

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get namespace from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.Namespace, error) {
	nsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ns := &corev1.Namespace{}
	if err = json.Unmarshal(nsJson, ns); err != nil {
		return nil, err
	}

	return h.GetByName(ns.Name)
}

// GetFromFile get namespace from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.Namespace, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get namespace by name.
func (h *Handler) GetByName(name string) (*corev1.Namespace, error) {
	return h.clientset.CoreV1().Namespaces().Get(h.ctx, name, h.Options.GetOptions)
}

// Get get namespace by name, alias to "GetByName".
func (h *Handler) Get(name string) (*corev1.Namespace, error) {
	return h.GetByName(name)
}
