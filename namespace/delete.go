package namespace

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete namespace from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	nsJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	ns := &corev1.Namespace{}
	err = json.Unmarshal(nsJson, ns)
	if err != nil {
		return err
	}

	return h.DeleteByName(ns.Name)
}

// DeleteFromFile delete namespace from yaml file.
func (h *Handler) DeleteFromFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete namespace by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().Namespaces().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete namespace by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
