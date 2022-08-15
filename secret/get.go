package secret

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets secret from type string, []byte, *corev1.Secret,
// corev1.Secret, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a secret from file path.
func (h *Handler) Get(obj interface{}) (*corev1.Secret, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *corev1.Secret:
		return h.GetFromObject(val)
	case corev1.Secret:
		return h.GetFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.GetFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.GetFromObject(val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	default:
		return nil, ErrInvalidGetType
	}
}

// GetByName gets secret by name.
func (h *Handler) GetByName(name string) (*corev1.Secret, error) {
	return h.clientset.CoreV1().Secrets(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets secret from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.Secret, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets secret from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.Secret, error) {
	secretJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{}
	if err = json.Unmarshal(secretJson, secret); err != nil {
		return nil, err
	}
	return h.getSecret(secret)
}

// GetFromObject gets secret from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*corev1.Secret, error) {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Secret")
	}
	return h.getSecret(secret)
}

// GetFromUnstructured gets secret from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), secret)
	if err != nil {
		return nil, err
	}
	return h.getSecret(secret)
}

// GetFromMap gets secret from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, secret)
	if err != nil {
		return nil, err
	}
	return h.getSecret(secret)
}

// getSecret
// It's necessary to get a new secret resource from a old secret resource,
// because old secret usually don't have secret.Status field.
func (h *Handler) getSecret(secret *corev1.Secret) (*corev1.Secret, error) {
	var namespace string
	if len(secret.Namespace) != 0 {
		namespace = secret.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().Secrets(namespace).Get(h.ctx, secret.Name, h.Options.GetOptions)
}
