package daemonset

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete daemonset from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	dsJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	daemonset := &appsv1.DaemonSet{}
	err = json.Unmarshal(dsJson, daemonset)
	if err != nil {
		return err
	}

	var namespace string
	if len(daemonset.Namespace) != 0 {
		namespace = daemonset.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(daemonset.Name)
}

// DeleteFromFile delete daemonset from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete daemonset by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.AppsV1().DaemonSets(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete daemonset by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
