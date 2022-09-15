package configmap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets configmap from type string, []byte, *corev1.ConfigMap,
// corev1.ConfigMap, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a configmap from file path.
func (h *Handler) Get(obj interface{}) (*corev1.ConfigMap, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *corev1.ConfigMap:
		return h.GetFromObject(val)
	case corev1.ConfigMap:
		return h.GetFromObject(&val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	case runtime.Object:
		return h.GetFromObject(val)
	default:
		return nil, ErrInvalidGetType
	}
}

// GetByName gets configmap by name.
func (h *Handler) GetByName(name string) (*corev1.ConfigMap, error) {
	return h.clientset.CoreV1().ConfigMaps(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets configmap from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.ConfigMap, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets configmap from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.ConfigMap, error) {
	cmJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cm := &corev1.ConfigMap{}
	if err = json.Unmarshal(cmJson, cm); err != nil {
		return nil, err
	}
	return h.getConfigmap(cm)
}

// GetFromObject gets configmap from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*corev1.ConfigMap, error) {
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.ConfigMap")
	}
	return h.getConfigmap(cm)
}

// GetFromUnstructured gets configmap from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cm)
	if err != nil {
		return nil, err
	}
	return h.getConfigmap(cm)
}

// GetFromMap gets configmap from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cm)
	if err != nil {
		return nil, err
	}
	return h.getConfigmap(cm)
}

// getConfigmap
// It's necessary to get a new configmap resource from a old configmap resource,
// because old configmap usually don't have configmap.Status field.
func (h *Handler) getConfigmap(cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	var namespace string
	if len(cm.Namespace) != 0 {
		namespace = cm.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().ConfigMaps(namespace).Get(h.ctx, cm.Name, h.Options.GetOptions)
}
