package persistentvolume

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create persistentvolume from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*corev1.PersistentVolume, error) {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, pv)
	if err != nil {
		return nil, err
	}

	return h.clientset.CoreV1().PersistentVolumes().Create(h.ctx, pv, h.Options.CreateOptions)
}

// CreateFromBytes create persistentvolume from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.PersistentVolume, error) {
	pvJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pv := &corev1.PersistentVolume{}
	if err = json.Unmarshal(pvJson, pv); err != nil {
		return nil, err
	}

	return h.clientset.CoreV1().PersistentVolumes().Create(h.ctx, pv, h.Options.CreateOptions)
}

// CreateFromFile create persistentvolume from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.PersistentVolume, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create persistentvolume from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*corev1.PersistentVolume, error) {
	return h.CreateFromFile(filename)
}
