package persistentvolume

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply persistentvolume from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*corev1.PersistentVolume, error) {
	pv := &corev1.PersistentVolume{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, pv)
	if err != nil {
		return nil, err
	}

	pv, err = h.clientset.CoreV1().PersistentVolumes().Create(h.ctx, pv, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		pv, err = h.clientset.CoreV1().PersistentVolumes().Update(h.ctx, pv, h.Options.UpdateOptions)
	}
	return pv, err
}

// ApplyFromBytes apply persistentvolume from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (pv *corev1.PersistentVolume, err error) {
	pv, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		pv, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply persistentvolume from yaml file.
func (h *Handler) ApplyFromFile(filename string) (pv *corev1.PersistentVolume, err error) {
	pv, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		pv, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply persistentvolume from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*corev1.PersistentVolume, error) {
	return h.ApplyFromFile(filename)
}
