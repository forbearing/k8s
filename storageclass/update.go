package storageclass

import (
	"encoding/json"
	"io/ioutil"

	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update storageclass from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*storagev1.StorageClass, error) {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, sc)
	if err != nil {
		return nil, err
	}

	return h.clientset.StorageV1().StorageClasses().Update(h.ctx, sc, h.Options.UpdateOptions)
}

// UpdateFromBytes update storageclass from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*storagev1.StorageClass, error) {
	scJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sc := &storagev1.StorageClass{}
	err = json.Unmarshal(scJson, sc)
	if err != nil {
		return nil, err
	}

	return h.clientset.StorageV1().StorageClasses().Update(h.ctx, sc, h.Options.UpdateOptions)
}

// UpdateFromFile update storageclass from yaml file
func (h *Handler) UpdateFromFile(filename string) (*storagev1.StorageClass, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update storageclass from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*storagev1.StorageClass, error) {
	return h.UpdateFromFile(filename)
}
