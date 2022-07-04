package replicationcontroller

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update replicationcontroller from map[string]interface{}
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*corev1.ReplicationController, error) {
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

	return h.clientset.CoreV1().ReplicationControllers(namespace).Update(h.ctx, rc, h.Options.UpdateOptions)
}

// UpdateFromBytes update replicationcontroller from bytes
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.ReplicationController, error) {
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

	return h.clientset.CoreV1().ReplicationControllers(namespace).Update(h.ctx, rc, h.Options.UpdateOptions)
}

// UpdateFromFile update replicationcontroller from yaml file
func (h *Handler) UpdateFromFile(filename string) (*corev1.ReplicationController, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update replicationcontroller from yaml file, alias to "UpdateFromFile"
func (h *Handler) Update(filename string) (*corev1.ReplicationController, error) {
	return h.UpdateFromFile(filename)
}
