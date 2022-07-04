package namespace

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create namespace from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*corev1.Namespace, error) {
	namespace := &corev1.Namespace{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, namespace)
	if err != nil {
		return nil, err
	}

	return h.clientset.CoreV1().Namespaces().Create(h.ctx, namespace, h.Options.CreateOptions)
}

// CreateFromBytes create namespace from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.Namespace, error) {
	nsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ns := &corev1.Namespace{}
	err = json.Unmarshal(nsJson, ns)
	if err != nil {
		return nil, err
	}

	return h.clientset.CoreV1().Namespaces().Create(h.ctx, ns, h.Options.CreateOptions)
}

// CreateFromFile create namespace from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.Namespace, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create namespace from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*corev1.Namespace, error) {
	return h.CreateFromFile(filename)
}
