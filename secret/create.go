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

// Create creates secret from type string, []byte, *corev1.Secret,
// corev1.Secret, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*corev1.Secret, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *corev1.Secret:
		return h.CreateFromObject(val)
	case corev1.Secret:
		return h.CreateFromObject(&val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates secret from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.Secret, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates secret from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.Secret, error) {
	secretJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	secret := &corev1.Secret{}
	if err = json.Unmarshal(secretJson, secret); err != nil {
		return nil, err
	}
	return h.createSecret(secret)
}

// CreateFromObject creates secret from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*corev1.Secret, error) {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Secret")
	}
	return h.createSecret(secret)
}

// CreateFromUnstructured creates secret from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), secret)
	if err != nil {
		return nil, err
	}
	return h.createSecret(secret)
}

// CreateFromMap creates secret from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*corev1.Secret, error) {
	secret := &corev1.Secret{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, secret)
	if err != nil {
		return nil, err
	}
	return h.createSecret(secret)
}

// createSecret
func (h *Handler) createSecret(secret *corev1.Secret) (*corev1.Secret, error) {
	var namespace string
	if len(secret.Namespace) != 0 {
		namespace = secret.Namespace
	} else {
		namespace = h.namespace
	}
	secret.ResourceVersion = ""
	secret.UID = ""
	return h.clientset.CoreV1().Secrets(namespace).Create(h.ctx, secret, h.Options.CreateOptions)
}
