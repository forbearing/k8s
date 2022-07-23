package replicaset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates replicaset from type string, []byte, *appsv1.ReplicaSet,
// appsv1.ReplicaSet, runtime.Object or map[string]interface{}.
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
	case runtime.Object:
		return h.CreateFromObject(val)
	case map[string]interface{}:
		return h.CreateFromUnstructured(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates replicaset from yaml file.
func (h *Handler) CreateFromFile(filename string) (*appsv1.ReplicaSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates replicaset from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*appsv1.ReplicaSet, error) {
	rsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rs := &appsv1.ReplicaSet{}
	err = json.Unmarshal(rsJson, rs)
	if err != nil {
		return nil, err
	}
	return h.createReplicaset(rs)
}

// CreateFromObject creates replicaset from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*appsv1.ReplicaSet, error) {
	rs, ok := obj.(*appsv1.ReplicaSet)
	if !ok {
		return nil, fmt.Errorf("object is not *appsv1.ReplicaSet")
	}
	return h.createReplicaset(rs)
}

// CreateFromUnstructured creates replicaset from map[string]interface{}.
func (h *Handler) CreateFromUnstructured(u map[string]interface{}) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rs)
	if err != nil {
		return nil, err
	}
	return h.createReplicaset(rs)
}

// createReplicaset
func (h *Handler) createReplicaset(rs *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error) {
	var namespace string
	if len(rs.Namespace) != 0 {
		namespace = rs.Namespace
	} else {
		namespace = h.namespace
	}
	rs.ResourceVersion = ""
	rs.UID = ""
	return h.clientset.AppsV1().ReplicaSets(namespace).Create(h.ctx, rs, h.Options.CreateOptions)
}
