package replicationcontroller

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

// Scale set replicationcontroller replicas from type string, []byte, *corev1.ReplicationController,
// corev1.ReplicationController, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call ScaleByName instead of ScaleFromFile.
// You should always explicitly call ScaleFromFile to set replicationcontroller replicas from file path.
func (h *Handler) Scale(obj interface{}, replicas int32) (*corev1.ReplicationController, error) {
	switch val := obj.(type) {
	case string:
		return h.ScaleByName(val, replicas)
	case []byte:
		return h.ScaleFromBytes(val, replicas)
	case *corev1.ReplicationController:
		return h.ScaleFromObject(val, replicas)
	case corev1.ReplicationController:
		return h.ScaleFromObject(&val, replicas)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.ScaleFromUnstructured(val.(*unstructured.Unstructured), replicas)
		}
		return h.ScaleFromObject(val, replicas)
	case *unstructured.Unstructured:
		return h.ScaleFromUnstructured(val, replicas)
	case unstructured.Unstructured:
		return h.ScaleFromUnstructured(&val, replicas)
	case map[string]interface{}:
		return h.ScaleFromMap(val, replicas)
	default:
		return nil, ErrInvalidScaleType
	}
}

// ScaleByName scale replicationcontroller by name.
func (h *Handler) ScaleByName(name string, replicas int32) (*corev1.ReplicationController, error) {
	rc, err := h.Get(name)
	if err != nil {
		return nil, err
	}
	copiedRC := rc.DeepCopy()
	if copiedRC.Spec.Replicas != nil {
		copiedRC.Spec.Replicas = &replicas
	}
	return h.Update(copiedRC)
}

// ScaleFromFile scale replicationcontroller from yaml file.
func (h *Handler) ScaleFromFile(filename string, replicas int32) (*corev1.ReplicationController, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.ScaleFromBytes(data, replicas)
}

// ScaleFromBytes scale replicationcontroller from bytes.
func (h *Handler) ScaleFromBytes(data []byte, replicas int32) (*corev1.ReplicationController, error) {
	rcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rc := &corev1.ReplicationController{}
	if err = json.Unmarshal(rcJson, rc); err != nil {
		return nil, err
	}
	return h.ScaleByName(rc.Name, replicas)
}

// ScaleFromObject scale replicationcontroller from runtime.Object.
func (h *Handler) ScaleFromObject(obj runtime.Object, replicas int32) (*corev1.ReplicationController, error) {
	rc, ok := obj.(*corev1.ReplicationController)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.ReplicationController")
	}
	return h.ScaleByName(rc.Name, replicas)
}

// ScaleFromUnstructured scale replicationcontroller from *unstructured.Unstructured.
func (h *Handler) ScaleFromUnstructured(u *unstructured.Unstructured, replicas int32) (*corev1.ReplicationController, error) {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rc)
	if err != nil {
		return nil, err
	}
	return h.ScaleByName(rc.Name, replicas)
}

// ScaleFromMap scale replicationcontroller from map[string]interface{}.
func (h *Handler) ScaleFromMap(u map[string]interface{}, replicas int32) (*corev1.ReplicationController, error) {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rc)
	if err != nil {
		return nil, err
	}
	return h.ScaleByName(rc.Name, replicas)
}
