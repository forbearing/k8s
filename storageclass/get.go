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

// Get gets storageclass from type string, []byte, *storagev1.StorageClass,
// storagev1.StorageClass, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a storageclass from file path.
func (h *Handler) Get(obj interface{}) (*storagev1.StorageClass, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *storagev1.StorageClass:
		return h.GetFromObject(val)
	case storagev1.StorageClass:
		return h.GetFromObject(&val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	case runtime.Object:
		return h.GetFromObject(val)
	default:
		return nil, ErrInvalidGetType
	}
}

// GetByName gets storageclass by name.
func (h *Handler) GetByName(name string) (*storagev1.StorageClass, error) {
	return h.clientset.StorageV1().StorageClasses().Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets storageclass from yaml file.
func (h *Handler) GetFromFile(filename string) (*storagev1.StorageClass, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets storageclass from bytes.
func (h *Handler) GetFromBytes(data []byte) (*storagev1.StorageClass, error) {
	scJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sc := &storagev1.StorageClass{}
	if err = json.Unmarshal(scJson, sc); err != nil {
		return nil, err
	}
	return h.getSC(sc)
}

// GetFromObject gets storageclass from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*storagev1.StorageClass, error) {
	sc, ok := obj.(*storagev1.StorageClass)
	if !ok {
		return nil, fmt.Errorf("object type is not *storagev1.StorageClass")
	}
	return h.getSC(sc)
}

// GetFromUnstructured gets storageclass from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*storagev1.StorageClass, error) {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sc)
	if err != nil {
		return nil, err
	}
	return h.getSC(sc)
}

// GetFromMap gets storageclass from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*storagev1.StorageClass, error) {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sc)
	if err != nil {
		return nil, err
	}
	return h.getSC(sc)
}

// getSC
// It's necessary to get a new storageclass resource from a old storageclass resource,
// because old storageclass usually don't have storageclass.Status field.
func (h *Handler) getSC(sc *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	return h.clientset.StorageV1().StorageClasses().Get(h.ctx, sc.Name, h.Options.GetOptions)
}
