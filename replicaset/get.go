package replicaset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets replicaset from type string, []byte, *appsv1.ReplicaSet,
// appsv1.ReplicaSet, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a replicaset from file path.
func (h *Handler) Get(obj interface{}) (*appsv1.ReplicaSet, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *appsv1.ReplicaSet:
		return h.GetFromObject(val)
	case appsv1.ReplicaSet:
		return h.GetFromObject(&val)
	case runtime.Object:
		return h.GetFromObject(val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets replicaset by name.
func (h *Handler) GetByName(name string) (*appsv1.ReplicaSet, error) {
	return h.clientset.AppsV1().ReplicaSets(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets replicaset from yaml file.
func (h *Handler) GetFromFile(filename string) (*appsv1.ReplicaSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets replicaset from bytes.
func (h *Handler) GetFromBytes(data []byte) (*appsv1.ReplicaSet, error) {
	rsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rs := &appsv1.ReplicaSet{}
	err = json.Unmarshal(rsJson, rs)
	if err != nil {
		return nil, err
	}
	return h.getReplicaset(rs)
}

// GetFromObject gets replicaset from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*appsv1.ReplicaSet, error) {
	rs, ok := obj.(*appsv1.ReplicaSet)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.ReplicaSet")
	}
	return h.getReplicaset(rs)
}

// GetFromUnstructured gets replicaset from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rs)
	if err != nil {
		return nil, err
	}
	return h.getReplicaset(rs)
}

// GetFromMap gets replicaset from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rs)
	if err != nil {
		return nil, err
	}
	return h.getReplicaset(rs)
}

// getReplicaset
// It's necessary to get a new replicaset resource from a old replicaset resource,
// because old replicaset usually don't have replicaset.Status field.
func (h *Handler) getReplicaset(rs *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error) {
	var namespace string
	if len(rs.Namespace) != 0 {
		namespace = rs.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.AppsV1().ReplicaSets(namespace).Get(h.ctx, rs.Name, h.Options.GetOptions)
}
