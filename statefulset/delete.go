package statefulset

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete statefulset from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	stsJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	sts := &appsv1.StatefulSet{}
	err = json.Unmarshal(stsJson, sts)
	if err != nil {
		return err
	}

	var namespace string
	if len(sts.Namespace) != 0 {
		namespace = sts.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(sts.Name)
}

// DeleteFromFile delete statefulset from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete statefulset by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.AppsV1().StatefulSets(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete statefulset by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
