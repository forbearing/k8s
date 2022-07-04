package replicaset

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete replicaset from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	dsJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	replicaset := &appsv1.ReplicaSet{}
	err = json.Unmarshal(dsJson, replicaset)
	if err != nil {
		return err
	}

	var namespace string
	if len(replicaset.Namespace) != 0 {
		namespace = replicaset.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(replicaset.Name)
}

// DeleteFromFile delete replicaset from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete replicaset by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.AppsV1().ReplicaSets(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete replicaset by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
