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

// Create creates replicaset from type string, []byte, *appsv1.ReplicaSet,
// appsv1.ReplicaSet, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*appsv1.ReplicaSet, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *appsv1.ReplicaSet:
		return h.CreateFromObject(val)
	case appsv1.ReplicaSet:
		return h.CreateFromObject(&val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.CreateFromObject(val)
	default:
		return nil, ErrInvalidCreateType
	}
}

// CreateFromFile creates replicaset from yaml or json file.
func (h *Handler) CreateFromFile(filename string) (*appsv1.ReplicaSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates replicaset from bytes data.
func (h *Handler) CreateFromBytes(data []byte) (*appsv1.ReplicaSet, error) {
	rsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rs := &appsv1.ReplicaSet{}
	if err = json.Unmarshal(rsJson, rs); err != nil {
		return nil, err
	}
	return h.createReplicaset(rs)
}

// CreateFromObject creates replicaset from metav1.Object or runtime.Object.
func (h *Handler) CreateFromObject(obj interface{}) (*appsv1.ReplicaSet, error) {
	rs, ok := obj.(*appsv1.ReplicaSet)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.ReplicaSet")
	}
	return h.createReplicaset(rs)
}

// CreateFromUnstructured creates replicaset from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rs)
	if err != nil {
		return nil, err
	}
	return h.createReplicaset(rs)
}

// CreateFromMap creates replicaset from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rs)
	if err != nil {
		return nil, err
	}
	return h.createReplicaset(rs)
}

// createReplicaset
func (h *Handler) createReplicaset(rs *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error) {
	namespace := rs.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	rs.ResourceVersion = ""
	rs.UID = ""
	return h.clientset.AppsV1().ReplicaSets(namespace).Create(h.ctx, rs, h.Options.CreateOptions)
}
