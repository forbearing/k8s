package namespace

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

// Delete deletes namespace from type string, []byte, *corev1.Namespace,
// corev1.Namespace, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a namespace from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *corev1.Namespace:
		return h.DeleteFromObject(val)
	case corev1.Namespace:
		return h.DeleteFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.DeleteFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.DeleteFromObject(val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	default:
		return ERR_TYPE_DELETE
	}
}

// DeleteByName deletes namespace by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().Namespaces().Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes namespace from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes namespace from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	nsJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	ns := &corev1.Namespace{}
	if err = json.Unmarshal(nsJson, ns); err != nil {
		return err
	}
	return h.deleteNamespace(ns)
}

// DeleteFromObject deletes namespace from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	ns, ok := obj.(*corev1.Namespace)
	if !ok {
		return fmt.Errorf("object type is not *corev1.Namespace")
	}
	return h.deleteNamespace(ns)
}

// DeleteFromUnstructured deletes namespace from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	ns := &corev1.Namespace{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), ns)
	if err != nil {
		return err
	}
	return h.deleteNamespace(ns)
}

// DeleteFromMap deletes namespace from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	ns := &corev1.Namespace{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, ns)
	if err != nil {
		return err
	}
	return h.deleteNamespace(ns)
}

// deleteNamespace
func (h *Handler) deleteNamespace(ns *corev1.Namespace) error {
	return h.clientset.CoreV1().Namespaces().Delete(h.ctx, ns.Name, h.Options.DeleteOptions)
}
