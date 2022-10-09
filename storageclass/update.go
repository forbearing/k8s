package storageclass

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates storageclass from type string, []byte, *storagev1.StorageClass,
// storagev1.StorageClass, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*storagev1.StorageClass, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *storagev1.StorageClass:
		return h.UpdateFromObject(val)
	case storagev1.StorageClass:
		return h.UpdateFromObject(&val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.UpdateFromObject(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates storageclass from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*storagev1.StorageClass, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates storageclass from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*storagev1.StorageClass, error) {
	scJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sc := &storagev1.StorageClass{}
	if err = json.Unmarshal(scJson, sc); err != nil {
		return nil, err
	}
	return h.updateSC(sc)
}

// UpdateFromObject updates storageclass from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*storagev1.StorageClass, error) {
	sc, ok := obj.(*storagev1.StorageClass)
	if !ok {
		return nil, fmt.Errorf("object type is not *storagev1.StorageClass")
	}
	return h.updateSC(sc)
}

// UpdateFromUnstructured updates storageclass from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*storagev1.StorageClass, error) {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sc)
	if err != nil {
		return nil, err
	}
	return h.updateSC(sc)
}

// UpdateFromMap updates storageclass from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*storagev1.StorageClass, error) {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sc)
	if err != nil {
		return nil, err
	}
	return h.updateSC(sc)
}

// updateSC
func (h *Handler) updateSC(sc *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	sc.ResourceVersion = ""
	sc.UID = ""
	return h.clientset.StorageV1().StorageClasses().Update(h.ctx, sc, h.Options.UpdateOptions)
}
