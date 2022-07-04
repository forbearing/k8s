package persistentvolume

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete persistentvolume from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	pvJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	pv := &corev1.PersistentVolume{}
	if err = json.Unmarshal(pvJson, pv); err != nil {
		return err
	}

	return h.DeleteByName(pv.Name)
}

// DeleteFromFile delete persistentvolume from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete persistentvolume by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().PersistentVolumes().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete persistentvolume by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
