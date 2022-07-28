package secret

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies secret from type string, []byte, *corev1.Secret,
// corev1.Secret, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*corev1.Secret, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *corev1.Secret:
		return h.ApplyFromObject(val)
	case corev1.Secret:
		return h.ApplyFromObject(&val)
	case runtime.Object:
		return h.ApplyFromObject(val)
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	default:
		return nil, ERR_TYPE_APPLY
	}
}

// ApplyFromFile applies secret from yaml file.
func (h *Handler) ApplyFromFile(filename string) (secret *corev1.Secret, err error) {
	secret, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if secret already exist, update it.
		secret, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply secret from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (secret *corev1.Secret, err error) {
	secret, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		secret, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies secret from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*corev1.Secret, error) {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Secret")
	}
	return h.applySecret(secret)
}

// ApplyFromUnstructured applies secret from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), secret)
	if err != nil {
		return nil, err
	}
	return h.applySecret(secret)
}

// ApplyFromMap applies secret from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, secret)
	if err != nil {
		return nil, err
	}
	return h.applySecret(secret)
}

// applySecret
func (h *Handler) applySecret(secret *corev1.Secret) (*corev1.Secret, error) {
	_, err := h.createSecret(secret)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateSecret(secret)
	}
	return secret, err
}
