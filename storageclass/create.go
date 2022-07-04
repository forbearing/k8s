package storageclass

import (
	"encoding/json"
	"io/ioutil"

	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create storageclass from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*storagev1.StorageClass, error) {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, sc)
	if err != nil {
		return nil, err
	}

	return h.clientset.StorageV1().StorageClasses().Create(h.ctx, sc, h.Options.CreateOptions)
}

// CreateFromBytes create storageclass from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*storagev1.StorageClass, error) {
	scJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sc := &storagev1.StorageClass{}
	if err = json.Unmarshal(scJson, sc); err != nil {
		return nil, err
	}

	return h.clientset.StorageV1().StorageClasses().Create(h.ctx, sc, h.Options.CreateOptions)
}

// CreateFromFile create storageclass from yaml file.
func (h *Handler) CreateFromFile(filename string) (*storagev1.StorageClass, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create storageclass from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*storagev1.StorageClass, error) {
	return h.CreateFromFile(filename)
}
