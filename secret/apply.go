package secret

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply secret from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, secret)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(secret.Namespace) != 0 {
		namespace = secret.Namespace
	} else {
		namespace = h.namespace
	}

	secret, err = h.clientset.CoreV1().Secrets(namespace).Create(h.ctx, secret, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		secret, err = h.clientset.CoreV1().Secrets(namespace).Update(h.ctx, secret, h.Options.UpdateOptions)
	}
	return secret, err
}

// ApplyFromBytes apply secret from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (secret *corev1.Secret, err error) {
	secret, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		secret, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply secret from yaml file.
func (h *Handler) ApplyFromFile(filename string) (secret *corev1.Secret, err error) {
	secret, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		secret, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply secret from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*corev1.Secret, error) {
	return h.ApplyFromFile(filename)
}
