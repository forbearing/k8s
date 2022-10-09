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

// Create creates statefulset from type string, []byte, *appsv1.StatefulSet,
// appsv1.StatefulSet, metav1.Object, runtime.Object, *unstructured.Unstructured,
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

// CreateFromFile creates statefulset from yaml or json file.
func (h *Handler) CreateFromFile(filename string) (*appsv1.StatefulSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates statefulset from bytes data.
func (h *Handler) CreateFromBytes(data []byte) (*appsv1.StatefulSet, error) {
	stsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sts := &appsv1.StatefulSet{}
	if err = json.Unmarshal(stsJson, sts); err != nil {
		return nil, err
	}
	return h.createStatefulset(sts)
}

// CreateFromObject creates statefulset from metav1.Object or runtime.Object.
func (h *Handler) CreateFromObject(obj interface{}) (*appsv1.StatefulSet, error) {
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
	namespace := sts.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	sts.ResourceVersion = ""
	sts.UID = ""
	return h.clientset.AppsV1().StatefulSets(namespace).Create(h.ctx, sts, h.Options.CreateOptions)
}
