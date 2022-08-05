package configmap

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

// Update updates configmap from type string, []byte, *corev1.ConfigMap,
// corev1.ConfigMap, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*corev1.ConfigMap, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *corev1.ConfigMap:
		return h.UpdateFromObject(val)
	case corev1.ConfigMap:
		return h.UpdateFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.UpdateFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.UpdateFromObject(val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates configmap from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.ConfigMap, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates configmap from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.ConfigMap, error) {
	cmJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cm := &corev1.ConfigMap{}
	if err = json.Unmarshal(cmJson, cm); err != nil {
		return nil, err
	}
	return h.updateConfigmap(cm)
}

// UpdateFromObject updates configmap from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*corev1.ConfigMap, error) {
	cm, ok := obj.(*corev1.ConfigMap)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.ConfigMap")
	}
	return h.updateConfigmap(cm)
}

// UpdateFromUnstructured updates configmap from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cm)
	if err != nil {
		return nil, err
	}
	return h.updateConfigmap(cm)
}

// UpdateFromMap updates configmap from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cm)
	if err != nil {
		return nil, err
	}
	return h.updateConfigmap(cm)
}

// updateConfigmap
func (h *Handler) updateConfigmap(cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	var namespace string
	if len(cm.Namespace) != 0 {
		namespace = cm.Namespace
	} else {
		namespace = h.namespace
	}
	cm.ResourceVersion = ""
	cm.UID = ""
	return h.clientset.CoreV1().ConfigMaps(namespace).Update(h.ctx, cm, h.Options.UpdateOptions)
}
