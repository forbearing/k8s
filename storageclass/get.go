package storageclass

import (
	"encoding/json"
	"io/ioutil"

	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get storageclass from bytes.
func (h *Handler) GetFromBytes(data []byte) (*storagev1.StorageClass, error) {
	scJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sc := &storagev1.StorageClass{}
	err = json.Unmarshal(scJson, sc)
	if err != nil {
		return nil, err
	}

	return h.GetByName(sc.Name)
}

// GetFromFile get storageclass from yaml file.
func (h *Handler) GetFromFile(filename string) (*storagev1.StorageClass, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get storageclass by name
func (h *Handler) GetByName(name string) (*storagev1.StorageClass, error) {
	return h.clientset.StorageV1().StorageClasses().Get(h.ctx, name, h.Options.GetOptions)
}

// Get get storageclass by name, alias to "GetByName.
func (h *Handler) Get(name string) (*storagev1.StorageClass, error) {
	return h.GetByName(name)
}
