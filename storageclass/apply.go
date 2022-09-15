package storageclass

import (
	"fmt"

	storagev1 "k8s.io/api/storage/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies storageclass from type string, []byte, *storagev1.StorageClass,
// storagev1.StorageClass, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*storagev1.StorageClass, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *storagev1.StorageClass:
		return h.ApplyFromObject(val)
	case storagev1.StorageClass:
		return h.ApplyFromObject(&val)
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	case runtime.Object:
		return h.ApplyFromObject(val)
	default:
		return nil, ErrInvalidApplyType
	}
}

// ApplyFromFile applies storageclass from yaml file.
func (h *Handler) ApplyFromFile(filename string) (sc *storagev1.StorageClass, err error) {
	sc, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if storageclass already exist, update it.
		sc, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply storageclass from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (sc *storagev1.StorageClass, err error) {
	sc, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		sc, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies storageclass from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*storagev1.StorageClass, error) {
	sc, ok := obj.(*storagev1.StorageClass)
	if !ok {
		return nil, fmt.Errorf("object type is not *storagev1.StorageClass")
	}
	return h.applySC(sc)
}

// ApplyFromUnstructured applies storageclass from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*storagev1.StorageClass, error) {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sc)
	if err != nil {
		return nil, err
	}
	return h.applySC(sc)
}

// ApplyFromMap applies storageclass from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*storagev1.StorageClass, error) {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sc)
	if err != nil {
		return nil, err
	}
	return h.applySC(sc)
}

// applySC
func (h *Handler) applySC(sc *storagev1.StorageClass) (*storagev1.StorageClass, error) {
	_, err := h.createSC(sc)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateSC(sc)
	}
	return sc, err
}
