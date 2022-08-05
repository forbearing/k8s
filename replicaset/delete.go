package replicaset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Delete deletes replicaset from type string, []byte, *appsv1.ReplicaSet,
// appsv1.ReplicaSet, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a replicaset from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *appsv1.ReplicaSet:
		return h.DeleteFromObject(val)
	case appsv1.ReplicaSet:
		return h.DeleteFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.DeleteFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.DeleteFromObject(val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	default:
		return ErrInvalidDeleteType
	}
}

// DeleteByName deletes replicaset by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.AppsV1().ReplicaSets(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes replicaset from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes replicaset from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	rsJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	rs := &appsv1.ReplicaSet{}
	if err = json.Unmarshal(rsJson, rs); err != nil {
		return err
	}
	return h.deleteReplicaset(rs)
}

// DeleteFromObject deletes replicaset from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	rs, ok := obj.(*appsv1.ReplicaSet)
	if !ok {
		return fmt.Errorf("object type is not *appsv1.ReplicaSet")
	}
	return h.deleteReplicaset(rs)
}

// DeleteFromUnstructured deletes replicaset from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rs)
	if err != nil {
		return err
	}
	return h.deleteReplicaset(rs)
}

// DeleteFromMap deletes replicaset from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rs)
	if err != nil {
		return err
	}
	return h.deleteReplicaset(rs)
}

// deleteReplicaset
func (h *Handler) deleteReplicaset(rs *appsv1.ReplicaSet) error {
	var namespace string
	if len(rs.Namespace) != 0 {
		namespace = rs.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.AppsV1().ReplicaSets(namespace).Delete(h.ctx, rs.Name, h.Options.DeleteOptions)
}
