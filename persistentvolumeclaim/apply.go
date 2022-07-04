package persistentvolumeclaim

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply persistentvolumeclaim from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, pvc)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(pvc.Namespace) != 0 {
		namespace = pvc.Namespace
	} else {
		namespace = h.namespace
	}

	pvc, err = h.clientset.CoreV1().PersistentVolumeClaims(namespace).Create(h.ctx, pvc, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		pvc, err = h.clientset.CoreV1().PersistentVolumeClaims(namespace).Update(h.ctx, pvc, h.Options.UpdateOptions)
	}
	return pvc, err
}

// ApplyFromBytes apply persistentvolumeclaim from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (pvc *corev1.PersistentVolumeClaim, err error) {
	pvc, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		pvc, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply persistentvolumeclaim from yaml file.
func (h *Handler) ApplyFromFile(filename string) (pvc *corev1.PersistentVolumeClaim, err error) {
	pvc, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		pvc, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply persistentvolumeclaim from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*corev1.PersistentVolumeClaim, error) {
	return h.ApplyFromFile(filename)
}
