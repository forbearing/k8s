package replicationcontroller

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get replicationcontroller from bytes
func (h *Handler) GetFromBytes(data []byte) (*corev1.ReplicationController, error) {
	rcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rc := &corev1.ReplicationController{}
	if err = json.Unmarshal(rcJson, rc); err != nil {
		return nil, err
	}

	var namespace string
	if len(rc.Namespace) != 0 {
		namespace = rc.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).GetByName(rc.Name)
}

// GetFromFile get replicationcontroller from yaml file
func (h *Handler) GetFromFile(filename string) (*corev1.ReplicationController, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get replicationcontroller by name
func (h *Handler) GetByName(name string) (*corev1.ReplicationController, error) {
	return h.clientset.CoreV1().ReplicationControllers(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get replicationcontroller by name
func (h *Handler) Get(name string) (replicationcontroller *corev1.ReplicationController, err error) {
	return h.GetByName(name)
}
