package secret

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Delete deletes secret from type string, []byte, *corev1.Secret,
// corev1.Secret, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a secret from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *corev1.Secret:
		return h.DeleteFromObject(val)
	case corev1.Secret:
		return h.DeleteFromObject(&val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	case runtime.Object:
		return h.DeleteFromObject(val)
	default:
		return ErrInvalidDeleteType
	}
}

// DeleteByName deletes secret by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().Secrets(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes secret from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes secret from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	secretJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	secret := &corev1.Secret{}
	if err = json.Unmarshal(secretJson, secret); err != nil {
		return err
	}
	return h.deleteSecret(secret)
}

// DeleteFromObject deletes secret from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return fmt.Errorf("object type is not *corev1.Secret")
	}
	return h.deleteSecret(secret)
}

// DeleteFromUnstructured deletes secret from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	secret := &corev1.Secret{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), secret)
	if err != nil {
		return err
	}
	return h.deleteSecret(secret)
}

// DeleteFromMap deletes secret from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	secret := &corev1.Secret{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, secret)
	if err != nil {
		return err
	}
	return h.deleteSecret(secret)
}

// deleteSecret
func (h *Handler) deleteSecret(secret *corev1.Secret) error {
	var namespace string
	if len(secret.Namespace) != 0 {
		namespace = secret.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().Secrets(namespace).Delete(h.ctx, secret.Name, h.Options.DeleteOptions)
}
