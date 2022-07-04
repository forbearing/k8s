package storageclass

import (
	storagev1 "k8s.io/api/storage/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply storageclass from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*storagev1.StorageClass, error) {
	sc := &storagev1.StorageClass{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, sc)
	if err != nil {
		return nil, err
	}

	sc, err = h.clientset.StorageV1().StorageClasses().Create(h.ctx, sc, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		sc, err = h.clientset.StorageV1().StorageClasses().Update(h.ctx, sc, h.Options.UpdateOptions)
	}
	return sc, err
}

// ApplyFromBytes apply storageclass from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (sc *storagev1.StorageClass, err error) {
	sc, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		sc, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply storageclass from yaml file.
func (h *Handler) ApplyFromFile(filename string) (sc *storagev1.StorageClass, err error) {
	sc, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		sc, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply storageclass from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*storagev1.StorageClass, error) {
	return h.ApplyFromFile(filename)
}
