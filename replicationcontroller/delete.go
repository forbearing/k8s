package replicationcontroller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Delete deletes replicationcontroller from type string, []byte,
// *corev1.ReplicationController, corev1.ReplicationController, runtime.Object,
// *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a replicationcontroller from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *corev1.ReplicationController:
		return h.DeleteFromObject(val)
	case corev1.ReplicationController:
		return h.DeleteFromObject(&val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	case runtime.Object:
		return h.DeleteFromObject(val)
	default:
		return ErrInvalidDeleteType
	}
}

// DeleteByName deletes replicationcontroller by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().ReplicationControllers(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes replicationcontroller from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes replicationcontroller from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	rcJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	rc := &corev1.ReplicationController{}
	if err = json.Unmarshal(rcJson, rc); err != nil {
		return err
	}
	return h.deleteRC(rc)
}

// DeleteFromObject deletes replicationcontroller from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	rc, ok := obj.(*corev1.ReplicationController)
	if !ok {
		return fmt.Errorf("object type is not *corev1.ReplicationController")
	}
	return h.deleteRC(rc)
}

// DeleteFromUnstructured deletes replicationcontroller from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rc)
	if err != nil {
		return err
	}
	return h.deleteRC(rc)
}

// DeleteFromMap deletes replicationcontroller from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rc)
	if err != nil {
		return err
	}
	return h.deleteRC(rc)
}

// deleteRC
func (h *Handler) deleteRC(rc *corev1.ReplicationController) error {
	var namespace string
	if len(rc.Namespace) != 0 {
		namespace = rc.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().ReplicationControllers(namespace).Delete(h.ctx, rc.Name, h.Options.DeleteOptions)
}
