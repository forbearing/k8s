package replicationcontroller

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create replicationcontroller from map[string]interface{}
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*corev1.ReplicationController, error) {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, rc)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(rc.Namespace) != 0 {
		namespace = rc.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.CoreV1().ReplicationControllers(namespace).Create(h.ctx, rc, h.Options.CreateOptions)
}

// CreateFromBytes create replicationcontroller from bytes
func (h *Handler) CreateFromBytes(data []byte) (*corev1.ReplicationController, error) {
	rcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rc := &corev1.ReplicationController{}
	err = json.Unmarshal(rcJson, rc)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(rc.Namespace) != 0 {
		namespace = rc.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.CoreV1().ReplicationControllers(namespace).Create(h.ctx, rc, h.Options.CreateOptions)
}

// CreateFromFile create replicationcontroller from yaml file
func (h *Handler) CreateFromFile(filename string) (*corev1.ReplicationController, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create replicationcontroller from yaml file, alias to "CreateFromFile"
func (h *Handler) Create(filename string) (*corev1.ReplicationController, error) {
	return h.CreateFromFile(filename)
}
