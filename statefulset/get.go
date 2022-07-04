package statefulset

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get statefulset from bytes.
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

	var namespace string
	if len(sts.Namespace) != 0 {
		namespace = sts.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).GetByName(sts.Name)
}

// GetFromFile get statefulset from file.
func (h *Handler) GetFromFile(filename string) (*appsv1.StatefulSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get statefulset by name.
func (h *Handler) GetByName(name string) (*appsv1.StatefulSet, error) {
	return h.clientset.AppsV1().StatefulSets(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get statefulset by name, alias to "GetByName".
func (h *Handler) Get(name string) (*appsv1.StatefulSet, error) {
	return h.GetByName(name)
}
