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

// Update updates replicaset from type string, []byte, *appsv1.ReplicaSet,
// appsv1.ReplicaSet, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*appsv1.ReplicaSet, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *appsv1.ReplicaSet:
		return h.UpdateFromObject(val)
	case appsv1.ReplicaSet:
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

// UpdateFromFile updates replicaset from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*appsv1.ReplicaSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates replicaset from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*appsv1.ReplicaSet, error) {
	rsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rs := &appsv1.ReplicaSet{}
	if err = json.Unmarshal(rsJson, rs); err != nil {
		return nil, err
	}
	return h.updateReplicaset(rs)
}

// UpdateFromObject updates replicaset from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*appsv1.ReplicaSet, error) {
	rs, ok := obj.(*appsv1.ReplicaSet)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.ReplicaSet")
	}
	return h.updateReplicaset(rs)
}

// UpdateFromUnstructured updates replicaset from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rs)
	if err != nil {
		return nil, err
	}
	return h.updateReplicaset(rs)
}

// UpdateFromMap updates replicaset from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rs)
	if err != nil {
		return nil, err
	}
	return h.updateReplicaset(rs)
}

// updateReplicaset
func (h *Handler) updateReplicaset(rs *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error) {
	namespace := rs.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	rs.ResourceVersion = ""
	rs.UID = ""
	return h.clientset.AppsV1().ReplicaSets(namespace).Update(h.ctx, rs, h.Options.UpdateOptions)
}
