package persistentvolume

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update persistentvolume from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*corev1.PersistentVolume, error) {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, pv)
	if err != nil {
		return nil, err
	}

	return h.clientset.CoreV1().PersistentVolumes().Update(h.ctx, pv, h.Options.UpdateOptions)
}

// UpdateFromBytes update persistentvolume from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.PersistentVolume, error) {
	pvJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pv := &corev1.PersistentVolume{}
	err = json.Unmarshal(pvJson, pv)
	if err != nil {
		return nil, err
	}

	return h.clientset.CoreV1().PersistentVolumes().Update(h.ctx, pv, h.Options.UpdateOptions)
}

// UpdateFromFile update persistentvolume from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.PersistentVolume, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update persistentvolume from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*corev1.PersistentVolume, error) {
	return h.UpdateFromFile(filename)
}
