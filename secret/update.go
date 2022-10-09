package secret

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates secret from type string, []byte, *corev1.Secret,
// corev1.Secret, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*corev1.Secret, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *corev1.Secret:
		return h.UpdateFromObject(val)
	case corev1.Secret:
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

// UpdateFromFile updates secret from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.Secret, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates secret from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.Secret, error) {
	secretJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{}
	if err = json.Unmarshal(secretJson, secret); err != nil {
		return nil, err
	}
	return h.updateSecret(secret)
}

// UpdateFromObject updates secret from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*corev1.Secret, error) {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Secret")
	}
	return h.updateSecret(secret)
}

// UpdateFromUnstructured updates secret from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), secret)
	if err != nil {
		return nil, err
	}
	return h.updateSecret(secret)
}

// UpdateFromMap updates secret from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, secret)
	if err != nil {
		return nil, err
	}
	return h.updateSecret(secret)
}

// updateSecret
func (h *Handler) updateSecret(secret *corev1.Secret) (*corev1.Secret, error) {
	namespace := secret.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	secret.ResourceVersion = ""
	secret.UID = ""
	return h.clientset.CoreV1().Secrets(namespace).Update(h.ctx, secret, h.Options.UpdateOptions)
}
