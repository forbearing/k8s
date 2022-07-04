package persistentvolume

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get persistentvolume from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.PersistentVolume, error) {
	pvJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pv := &corev1.PersistentVolume{}
	err = json.Unmarshal(pvJson, pv)
	if err != nil {
		return nil, err
	}

	return h.GetByName(pv.Name)
}

// GetFromFile get persistentvolume from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.PersistentVolume, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get persistentvolume by name
func (h *Handler) GetByName(name string) (*corev1.PersistentVolume, error) {
	return h.clientset.CoreV1().PersistentVolumes().Get(h.ctx, name, h.Options.GetOptions)
}

// Get get persistentvolume by name, alias to "GetByName".
func (h *Handler) Get(name string) (*corev1.PersistentVolume, error) {
	return h.GetByName(name)
}
