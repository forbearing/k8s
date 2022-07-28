package statefulset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates statefulset from type string, []byte, *appsv1.StatefulSet,
// appsv1.StatefulSet, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*appsv1.StatefulSet, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *appsv1.StatefulSet:
		return h.CreateFromObject(val)
	case appsv1.StatefulSet:
		return h.CreateFromObject(&val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates statefulset from yaml file.
func (h *Handler) CreateFromFile(filename string) (*appsv1.StatefulSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates statefulset from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*appsv1.StatefulSet, error) {
	stsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sts := &appsv1.StatefulSet{}
	err = json.Unmarshal(stsJson, sts)
	if err != nil {
		return nil, err
	}
	return h.createStatefulset(sts)
}

// CreateFromObject creates statefulset from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*appsv1.StatefulSet, error) {
	sts, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		return nil, fmt.Errorf("object type is not *appsv1.StatefulSet")
	}
	return h.createStatefulset(sts)
}

// CreateFromUnstructured creates statefulset from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), sts)
	if err != nil {
		return nil, err
	}
	return h.createStatefulset(sts)
}

// CreateFromMap creates statefulset from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sts)
	if err != nil {
		return nil, err
	}
	return h.createStatefulset(sts)
}

// createStatefulset
func (h *Handler) createStatefulset(sts *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	var namespace string
	if len(sts.Namespace) != 0 {
		namespace = sts.Namespace
	} else {
		namespace = h.namespace
	}
	sts.ResourceVersion = ""
	sts.UID = ""
	return h.clientset.AppsV1().StatefulSets(namespace).Create(h.ctx, sts, h.Options.CreateOptions)
}
