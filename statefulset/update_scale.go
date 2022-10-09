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

// Scale set statefulset replicas from type string, []byte, *appsv1.StatefulSet,
// appsv1.StatefulSet, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call ScaleByName instead of ScaleFromFile.
// You should always explicitly call ScaleFromFile to set statefulset replicas from file path.
func (h *Handler) Scale(obj interface{}, replicas int32) (*appsv1.StatefulSet, error) {
	switch val := obj.(type) {
	case string:
		return h.ScaleByName(val, replicas)
	case []byte:
		return h.ScaleFromBytes(val, replicas)
	case *appsv1.StatefulSet:
		return h.ScaleFromObject(val, replicas)
	case appsv1.StatefulSet:
		return h.ScaleFromObject(&val, replicas)
	case *unstructured.Unstructured:
		return h.ScaleFromUnstructured(val, replicas)
	case unstructured.Unstructured:
		return h.ScaleFromUnstructured(&val, replicas)
	case map[string]interface{}:
		return h.ScaleFromMap(val, replicas)
	case metav1.Object, runtime.Object:
		return h.ScaleFromObject(val, replicas)
	default:
		return nil, ErrInvalidScaleType
	}
}

// ScaleByName scale statefulset by name.
func (h *Handler) ScaleByName(name string, replicas int32) (*appsv1.StatefulSet, error) {
	sts, err := h.Get(name)
	if err != nil {
		return nil, err
	}
	copiedSts := sts.DeepCopy()
	if copiedSts.Spec.Replicas != nil {
		copiedSts.Spec.Replicas = &replicas
	}
	return h.Update(copiedSts)
}

// ScaleFromFile scale statefulset from yaml or json file.
func (h *Handler) ScaleFromFile(filename string, replicas int32) (*appsv1.StatefulSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.ScaleFromBytes(data, replicas)
}

// ScaleFromBytes scale statefulset from bytes data.
func (h *Handler) ScaleFromBytes(data []byte, replicas int32) (*appsv1.StatefulSet, error) {
	stsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sts := &appsv1.StatefulSet{}
	if err = json.Unmarshal(stsJson, sts); err != nil {
		return nil, err
	}
	return h.ScaleByName(sts.Name, replicas)
}

// ScaleFromObject scale statefulset from metav1.Object or runtime.Object.
func (h *Handler) ScaleFromObject(obj interface{}, replicas int32) (*appsv1.StatefulSet, error) {
	sts, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.StatefulSet")
	}
	return h.ScaleByName(sts.Name, replicas)
}

// ScaleFromUnstructured scale statefulset from *unstructured.Unstructured.
func (h *Handler) ScaleFromUnstructured(u *unstructured.Unstructured, replicas int32) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sts)
	if err != nil {
		return nil, err
	}
	return h.ScaleByName(sts.Name, replicas)
}

// ScaleFromMap scale statefulset from map[string]interface{}.
func (h *Handler) ScaleFromMap(u map[string]interface{}, replicas int32) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sts)
	if err != nil {
		return nil, err
	}
	return h.ScaleByName(sts.Name, replicas)
}
