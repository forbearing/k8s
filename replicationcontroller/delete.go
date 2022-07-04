package replicationcontroller

import (
	"encoding/json"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete replicationcontroller from bytes
func (h *Handler) DeleteFromBytes(data []byte) error {
	rcJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	rc := &corev1.ReplicationController{}
	err = json.Unmarshal(rcJson, rc)
	if err != nil {
		return err
	}

	var namespace string
	if len(rc.Namespace) != 0 {
		namespace = rc.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(rc.Name)
}

// DeleteFromFile delete replicationcontroller from yaml file
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete replicationcontroller by name
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().ReplicationControllers(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete replicationcontroller by name, alias to "DeleteByName"
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
