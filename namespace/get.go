package namespace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets namespace from type string, []byte, *corev1.Namespace,
// corev1.Namespace, runtime.Object or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a namespace from file path.
func (h *Handler) Get(obj interface{}) (*corev1.Namespace, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *corev1.Namespace:
		return h.GetFromObject(val)
	case corev1.Namespace:
		return h.GetFromObject(&val)
	case map[string]interface{}:
		return h.GetFromUnstructured(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets namespace by name.
func (h *Handler) GetByName(name string) (*corev1.Namespace, error) {
	return h.clientset.CoreV1().Namespaces().Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets namespace from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.Namespace, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets namespace from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.Namespace, error) {
	nsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ns := &corev1.Namespace{}
	err = json.Unmarshal(nsJson, ns)
	if err != nil {
		return nil, err
	}
	return h.getNamespace(ns)
}

// GetFromObject gets namespace from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*corev1.Namespace, error) {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		return nil, fmt.Errorf("object is not *corev1.Namespace")
	}
	return h.getNamespace(ns)
}

// GetFromUnstructured gets namespace from map[string]interface{}.
func (h *Handler) GetFromUnstructured(u map[string]interface{}) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ns)
	if err != nil {
		return nil, err
	}
	return h.getNamespace(ns)
}

// getNamespace
// It's necessary to get a new namespace resource from a old namespace resource,
// because old namespace usually don't have namespace.Status field.
func (h *Handler) getNamespace(ns *corev1.Namespace) (*corev1.Namespace, error) {
	return h.clientset.CoreV1().Namespaces().Get(h.ctx, ns.Name, h.Options.GetOptions)
}
