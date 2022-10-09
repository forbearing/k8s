package replicaset

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

// Scale set replicaset replicas from type string, []byte, *appsv1.ReplicaSet,
// appsv1.ReplicaSet, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call ScaleByName instead of ScaleFromFile.
// You should always explicitly call ScaleFromFile to set replicaset replicas from file path.
func (h *Handler) Scale(obj interface{}, replicas int32) (*appsv1.ReplicaSet, error) {
	switch val := obj.(type) {
	case string:
		return h.ScaleByName(val, replicas)
	case []byte:
		return h.ScaleFromBytes(val, replicas)
	case *appsv1.ReplicaSet:
		return h.ScaleFromObject(val, replicas)
	case appsv1.ReplicaSet:
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

// ScaleByName scale replicaset by name.
func (h *Handler) ScaleByName(name string, replicas int32) (*appsv1.ReplicaSet, error) {
	rs, err := h.Get(name)
	if err != nil {
		return nil, err
	}
	copiedRS := rs.DeepCopy()
	if copiedRS.Spec.Replicas != nil {
		copiedRS.Spec.Replicas = &replicas
	}
	return h.Update(copiedRS)
}

// ScaleFromFile scale replicaset from yaml or json file.
func (h *Handler) ScaleFromFile(filename string, replicas int32) (*appsv1.ReplicaSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.ScaleFromBytes(data, replicas)
}

// ScaleFromBytes scale replicaset from bytes data.
func (h *Handler) ScaleFromBytes(data []byte, replicas int32) (*appsv1.ReplicaSet, error) {
	rsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rs := &appsv1.ReplicaSet{}
	if err = json.Unmarshal(rsJson, rs); err != nil {
		return nil, err
	}
	return h.ScaleByName(rs.Name, replicas)
}

// ScaleFromObject scale replicaset from metav1.Object or runtime.Object.
func (h *Handler) ScaleFromObject(obj interface{}, replicas int32) (*appsv1.ReplicaSet, error) {
	rs, ok := obj.(*appsv1.ReplicaSet)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.ReplicaSet")
	}
	return h.ScaleByName(rs.Name, replicas)
}

// ScaleFromUnstructured scale replicaset from *unstructured.Unstructured.
func (h *Handler) ScaleFromUnstructured(u *unstructured.Unstructured, replicas int32) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rs)
	if err != nil {
		return nil, err
	}
	return h.ScaleByName(rs.Name, replicas)
}

// ScaleFromMap scale replicaset from map[string]interface{}.
func (h *Handler) ScaleFromMap(u map[string]interface{}, replicas int32) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rs)
	if err != nil {
		return nil, err
	}
	return h.ScaleByName(rs.Name, replicas)
}
