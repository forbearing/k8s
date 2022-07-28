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

// Create creates storageclass from type string, []byte, *storagev1.StorageClass,
// storagev1.StorageClass, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*storagev1.StorageClass, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *storagev1.StorageClass:
		return h.CreateFromObject(val)
	case storagev1.StorageClass:
		return h.CreateFromObject(&val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates storageclass from yaml file.
func (h *Handler) CreateFromFile(filename string) (*storagev1.StorageClass, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates storageclass from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*storagev1.StorageClass, error) {
	scJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sc := &storagev1.StorageClass{}
	err = json.Unmarshal(scJson, sc)
	if err != nil {
		return nil, err
	}
	return h.createSC(sc)
}

// CreateFromObject creates storageclass from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*storagev1.StorageClass, error) {
	sc, ok := obj.(*storagev1.StorageClass)
	if !ok {
		return nil, fmt.Errorf("object type is not *storagev1.StorageClass")
	}
	return h.createSC(sc)
}

// CreateFromUnstructured creates storageclass from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*storagev1.StorageClass, error) {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sc)
	if err != nil {
		return nil, err
	}
	return h.createSC(sc)
}

// CreateFromMap creates storageclass from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*storagev1.StorageClass, error) {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sc)
	if err != nil {
		return nil, err
	}
	return h.createSC(sc)
}

// createSC
func (h *Handler) createSC(sc *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	sc.ResourceVersion = ""
	sc.UID = ""
	return h.clientset.StorageV1().StorageClasses().Create(h.ctx, sc, h.Options.CreateOptions)
}
