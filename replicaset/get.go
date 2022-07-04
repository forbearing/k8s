package replicaset

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get replicaset from bytes..
func (h *Handler) GetFromBytes(data []byte) (*appsv1.ReplicaSet, error) {
	dsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	replicaset := &appsv1.ReplicaSet{}
	if err = json.Unmarshal(dsJson, replicaset); err != nil {
		return nil, err
	}

	var namespace string
	if len(replicaset.Namespace) != 0 {
		namespace = replicaset.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).GetByName(replicaset.Name)
}

// GetFromFile get replicaset from yaml file.
func (h *Handler) GetFromFile(filename string) (*appsv1.ReplicaSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return h.GetFromBytes(data)
}

// GetByName get replicaset by name.
func (h *Handler) GetByName(name string) (*appsv1.ReplicaSet, error) {
	return h.clientset.AppsV1().ReplicaSets(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get replicaset by name, alias to "GetByName".
func (h *Handler) Get(name string) (*appsv1.ReplicaSet, error) {
	return h.GetByName(name)
}
