package replicaset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates replicaset from type string, []byte, *appsv1.ReplicaSet,
// appsv1.ReplicaSet, runtime.Object or map[string]interface{}.
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
	case runtime.Object:
		return h.UpdateFromObject(val)
	case map[string]interface{}:
		return h.UpdateFromUnstructured(val)
	default:
		return nil, ERR_TYPE_UPDATE
	}
}

// UpdateFromFile updates replicaset from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*appsv1.ReplicaSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates replicaset from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*appsv1.ReplicaSet, error) {
	rsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rs := &appsv1.ReplicaSet{}
	err = json.Unmarshal(rsJson, rs)
	if err != nil {
		return nil, err
	}
	return h.updateReplicaset(rs)
}

// UpdateFromObject updates replicaset from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*appsv1.ReplicaSet, error) {
	rs, ok := obj.(*appsv1.ReplicaSet)
	if !ok {
		return nil, fmt.Errorf("object is not *appsv1.ReplicaSet")
	}
	return h.updateReplicaset(rs)
}

// UpdateFromUnstructured updates replicaset from map[string]interface{}.
func (h *Handler) UpdateFromUnstructured(u map[string]interface{}) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rs)
	if err != nil {
		return nil, err
	}
	return h.updateReplicaset(rs)
}

// updateReplicaset
func (h *Handler) updateReplicaset(rs *appsv1.ReplicaSet) (*appsv1.ReplicaSet, error) {
	var namespace string
	if len(rs.Namespace) != 0 {
		namespace = rs.Namespace
	} else {
		namespace = h.namespace
	}
	rs.ResourceVersion = ""
	rs.UID = ""
	return h.clientset.AppsV1().ReplicaSets(namespace).Update(h.ctx, rs, h.Options.UpdateOptions)
}
