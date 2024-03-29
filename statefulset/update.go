package statefulset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates statefulset from type string, []byte, *appsv1.StatefulSet,
// appsv1.StatefulSet, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*appsv1.StatefulSet, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *appsv1.StatefulSet:
		return h.UpdateFromObject(val)
	case appsv1.StatefulSet:
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

// UpdateFromFile updates statefulset from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*appsv1.StatefulSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates statefulset from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*appsv1.StatefulSet, error) {
	stsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sts := &appsv1.StatefulSet{}
	if err = json.Unmarshal(stsJson, sts); err != nil {
		return nil, err
	}
	return h.updateStatefulset(sts)
}

// UpdateFromObject updates statefulset from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*appsv1.StatefulSet, error) {
	sts, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.StatefulSet")
	}
	return h.updateStatefulset(sts)
}

// UpdateFromUnstructured updates statefulset from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sts)
	if err != nil {
		return nil, err
	}
	return h.updateStatefulset(sts)
}

// UpdateFromMap updates statefulset from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sts)
	if err != nil {
		return nil, err
	}
	return h.updateStatefulset(sts)
}

// updateStatefulset
func (h *Handler) updateStatefulset(sts *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	namespace := sts.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	sts.ResourceVersion = ""
	sts.UID = ""
	return h.clientset.AppsV1().StatefulSets(namespace).Update(h.ctx, sts, h.Options.UpdateOptions)
}
