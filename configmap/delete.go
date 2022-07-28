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

// Delete deletes configmap from type string, []byte, *corev1.ConfigMap,
// corev1.ConfigMap, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a configmap from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *corev1.ConfigMap:
		return h.DeleteFromObject(val)
	case corev1.ConfigMap:
		return h.DeleteFromObject(&val)
	case runtime.Object:
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

// DeleteByName deletes configmap by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().ConfigMaps(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes configmap from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes configmap from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	cmJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	cm := &corev1.ConfigMap{}
	err = json.Unmarshal(cmJson, cm)
	if err != nil {
		return err
	}
	return h.deleteConfigmap(cm)
}

// DeleteFromObject deletes configmap from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return fmt.Errorf("object type is not *corev1.ConfigMap")
	}
	return h.deleteConfigmap(cm)
}

// DeleteFromUnstructured deletes configmap from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	cm := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cm)
	if err != nil {
		return err
	}
	return h.deleteConfigmap(cm)
}

// DeleteFromMap deletes configmap from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	cm := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cm)
	if err != nil {
		return err
	}
	return h.deleteConfigmap(cm)
}

// deleteConfigmap
func (h *Handler) deleteConfigmap(cm *corev1.ConfigMap) error {
	var namespace string
	if len(cm.Namespace) != 0 {
		namespace = cm.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().ConfigMaps(namespace).Delete(h.ctx, cm.Name, h.Options.DeleteOptions)
}
