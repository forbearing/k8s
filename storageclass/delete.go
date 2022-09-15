package storageclass

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Delete deletes storageclass from type string, []byte, *storagev1.StorageClass,
// storagev1.StorageClass, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a storageclass from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *storagev1.StorageClass:
		return h.DeleteFromObject(val)
	case storagev1.StorageClass:
		return h.DeleteFromObject(&val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	case runtime.Object:
		return h.DeleteFromObject(val)
	default:
		return ErrInvalidDeleteType
	}
}

// DeleteByName deletes storageclass by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.StorageV1().StorageClasses().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes storageclass from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes storageclass from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	scJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	sc := &storagev1.StorageClass{}
	if err = json.Unmarshal(scJson, sc); err != nil {
		return err
	}
	return h.deleteSC(sc)
}

// DeleteFromObject deletes storageclass from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	sc, ok := obj.(*storagev1.StorageClass)
	if !ok {
		return fmt.Errorf("object type is not *storagev1.StorageClass")
	}
	return h.deleteSC(sc)
}

// DeleteFromUnstructured deletes storageclass from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sc)
	if err != nil {
		return err
	}
	return h.deleteSC(sc)
}

// DeleteFromMap deletes storageclass from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sc)
	if err != nil {
		return err
	}
	return h.deleteSC(sc)
}

// deleteSC
func (h *Handler) deleteSC(sc *storagev1.StorageClass) error {
	return h.clientset.StorageV1().StorageClasses().Delete(h.ctx, sc.Name, h.Options.DeleteOptions)
}
