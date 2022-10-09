package replicationcontroller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets replicationcontroller from type string, []byte,
// *corev1.ReplicationController, corev1.ReplicationController, metav1.Object, runtime.Object,
// *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a replicationcontroller from file path.
func (h *Handler) Get(obj interface{}) (*corev1.ReplicationController, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *corev1.ReplicationController:
		return h.GetFromObject(val)
	case corev1.ReplicationController:
		return h.GetFromObject(&val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	case metav1.Object, runtime.Object:
		return h.GetFromObject(val)
	default:
		return nil, ErrInvalidGetType
	}
}

// GetByName gets replicationcontroller by name.
func (h *Handler) GetByName(name string) (*corev1.ReplicationController, error) {
	return h.clientset.CoreV1().ReplicationControllers(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets replicationcontroller from yaml or json file.
func (h *Handler) GetFromFile(filename string) (*corev1.ReplicationController, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets replicationcontroller from bytes data.
func (h *Handler) GetFromBytes(data []byte) (*corev1.ReplicationController, error) {
	rcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rc := &corev1.ReplicationController{}
	if err = json.Unmarshal(rcJson, rc); err != nil {
		return nil, err
	}
	return h.getRS(rc)
}

// GetFromObject gets replicationcontroller from metav1.Object or runtime.Object.
func (h *Handler) GetFromObject(obj interface{}) (*corev1.ReplicationController, error) {
	rc, ok := obj.(*corev1.ReplicationController)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.ReplicationController")
	}
	return h.getRS(rc)
}

// GetFromUnstructured gets replicationcontroller from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*corev1.ReplicationController, error) {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rc)
	if err != nil {
		return nil, err
	}
	return h.getRS(rc)
}

// GetFromMap gets replicationcontroller from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*corev1.ReplicationController, error) {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rc)
	if err != nil {
		return nil, err
	}
	return h.getRS(rc)
}

// getRS
// It's necessary to get a new replicationcontroller resource from a old replicationcontroller resource,
// because old replicationcontroller usually don't have replicationcontroller.Status field.
func (h *Handler) getRS(rc *corev1.ReplicationController) (*corev1.ReplicationController, error) {
	namespace := rc.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().ReplicationControllers(namespace).Get(h.ctx, rc.Name, h.Options.GetOptions)
}
