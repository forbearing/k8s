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

// Update updates replicationcontroller from type string, []byte,
// *corev1.ReplicationController, corev1.ReplicationController, runtime.Object,
// *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*corev1.ReplicationController, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *corev1.ReplicationController:
		return h.UpdateFromObject(val)
	case corev1.ReplicationController:
		return h.UpdateFromObject(&val)
	case runtime.Object:
		return h.UpdateFromObject(val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	default:
		return nil, ERR_TYPE_UPDATE
	}
}

// UpdateFromFile updates replicationcontroller from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.ReplicationController, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates replicationcontroller from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.ReplicationController, error) {
	rcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rc := &corev1.ReplicationController{}
	if err = json.Unmarshal(rcJson, rc); err != nil {
		return nil, err
	}
	return h.updateRS(rc)
}

// UpdateFromObject updates replicationcontroller from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*corev1.ReplicationController, error) {
	rc, ok := obj.(*corev1.ReplicationController)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.ReplicationController")
	}
	return h.updateRS(rc)
}

// UpdateFromUnstructured updates replicationcontroller from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*corev1.ReplicationController, error) {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rc)
	if err != nil {
		return nil, err
	}
	return h.updateRS(rc)
}

// UpdateFromMap updates replicationcontroller from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*corev1.ReplicationController, error) {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rc)
	if err != nil {
		return nil, err
	}
	return h.updateRS(rc)
}

// updateRS
func (h *Handler) updateRS(rc *corev1.ReplicationController) (*corev1.ReplicationController, error) {
	var namespace string
	if len(rc.Namespace) != 0 {
		namespace = rc.Namespace
	} else {
		namespace = h.namespace
	}
	rc.ResourceVersion = ""
	rc.UID = ""
	return h.clientset.CoreV1().ReplicationControllers(namespace).Update(h.ctx, rc, h.Options.UpdateOptions)
}
