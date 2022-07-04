package storageclass

import (
	"encoding/json"
	"io/ioutil"

	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete storageclass from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	scJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	sc := &storagev1.StorageClass{}
	if err = json.Unmarshal(scJson, sc); err != nil {
		return err
	}

	return h.DeleteByName(sc.Name)
}

// DeleteFromFile delete storageclass from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete storageclass by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.StorageV1().StorageClasses().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete storageclass by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
