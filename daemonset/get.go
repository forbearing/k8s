package daemonset

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get daemonset from bytes.
func (h *Handler) GetFromBytes(data []byte) (*appsv1.DaemonSet, error) {
	dsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	daemonset := &appsv1.DaemonSet{}
	if err = json.Unmarshal(dsJson, daemonset); err != nil {
		return nil, err
	}

	var namespace string
	if len(daemonset.Namespace) != 0 {
		namespace = daemonset.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).GetByName(daemonset.Name)
}

// GetFromFile get daemonset from yaml file.
func (h *Handler) GetFromFile(filename string) (*appsv1.DaemonSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return h.GetFromBytes(data)
}

// GetByName get daemonset by name.
func (h *Handler) GetByName(name string) (*appsv1.DaemonSet, error) {
	return h.clientset.AppsV1().DaemonSets(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get daemonset by name, alias to "GetByName".
func (h *Handler) Get(name string) (*appsv1.DaemonSet, error) {
	return h.GetByName(name)
}
