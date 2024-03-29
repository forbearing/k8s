package namespace

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

// Update updates namespace from type string, []byte, *corev1.Namespace,
// corev1.Namespace, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*corev1.Namespace, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *corev1.Namespace:
		return h.UpdateFromObject(val)
	case corev1.Namespace:
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

// UpdateFromFile updates namespace from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.Namespace, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates namespace from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.Namespace, error) {
	nsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	ns := &corev1.Namespace{}
	if err = json.Unmarshal(nsJson, ns); err != nil {
		return nil, err
	}
	return h.updateNamespace(ns)
}

// UpdateFromObject updates namespace from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*corev1.Namespace, error) {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Namespace")
	}
	return h.updateNamespace(ns)
}

// UpdateFromUnstructured updates namespace from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ns)
	if err != nil {
		return nil, err
	}
	return h.updateNamespace(ns)
}

// UpdateFromMap updates namespace from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ns)
	if err != nil {
		return nil, err
	}
	return h.updateNamespace(ns)
}

// updateNamespace
func (h *Handler) updateNamespace(ns *corev1.Namespace) (*corev1.Namespace, error) {
	ns.ResourceVersion = ""
	ns.UID = ""
	return h.clientset.CoreV1().Namespaces().Update(h.ctx, ns, h.Options.UpdateOptions)
}
