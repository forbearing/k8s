package namespace

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply namespace from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*corev1.Namespace, error) {
	namespace := &corev1.Namespace{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, namespace)
	if err != nil {
		return nil, err
	}

	namespace, err = h.clientset.CoreV1().Namespaces().Create(h.ctx, namespace, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		namespace, err = h.clientset.CoreV1().Namespaces().Update(h.ctx, namespace, h.Options.UpdateOptions)
	}
	return namespace, err
}

// ApplyFromBytes apply namespace from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (namespace *corev1.Namespace, err error) {
	namespace, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		namespace, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply namespace from yaml file.
func (h *Handler) ApplyFromFile(filename string) (namespace *corev1.Namespace, err error) {
	namespace, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		namespace, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply namespace from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*corev1.Namespace, error) {
	return h.ApplyFromFile(filename)
}
