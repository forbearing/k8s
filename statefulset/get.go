package statefulset

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets statefulset from type string, []byte, *appsv1.StatefulSet,
// appsv1.StatefulSet, runtime.Object or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You need to always explicitly call GetFromFile to get a statefulset from file path.
func (h *Handler) Get(obj interface{}) (*appsv1.StatefulSet, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *appsv1.StatefulSet:
		return h.GetFromObject(val)
	case appsv1.StatefulSet:
		return h.GetFromObject(&val)
	case map[string]interface{}:
		return h.GetFromUnstructured(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets statefulset by name.
func (h *Handler) GetByName(name string) (*appsv1.StatefulSet, error) {
	return h.clientset.AppsV1().StatefulSets(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets statefulset from yaml file.
func (h *Handler) GetFromFile(filename string) (*appsv1.StatefulSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets statefulset from bytes.
func (h *Handler) GetFromBytes(data []byte) (*appsv1.StatefulSet, error) {
	stsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	sts := &appsv1.StatefulSet{}
	err = json.Unmarshal(stsJson, sts)
	if err != nil {
		return nil, err
	}
	return h.getStatefulset(sts)
}

// GetFromObject gets statefulset from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*appsv1.StatefulSet, error) {
	sts, ok := obj.(*appsv1.StatefulSet)
	if !ok {
		return nil, fmt.Errorf("object is not *appsv1.StatefulSet")
	}
	return h.getStatefulset(sts)
}

// GetFromUnstructured gets statefulset from map[string]interface{}.
func (h *Handler) GetFromUnstructured(u map[string]interface{}) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, sts)
	if err != nil {
		return nil, err
	}
	return h.getStatefulset(sts)
}

// getStatefulset
// It's necessary to get a new statefulset resource from a old statefulset resource,
// because old statefulset usually don't have statefulset.Status field.
func (h *Handler) getStatefulset(sts *appsv1.StatefulSet) (*appsv1.StatefulSet, error) {
	var namespace string
	if len(sts.Namespace) != 0 {
		namespace = sts.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.AppsV1().StatefulSets(namespace).Get(h.ctx, sts.Name, h.Options.GetOptions)
}
