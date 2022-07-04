package daemonset

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update daemonset from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*appsv1.DaemonSet, error) {
	daemonset := &appsv1.DaemonSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, daemonset)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(daemonset.Namespace) != 0 {
		namespace = daemonset.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.AppsV1().DaemonSets(namespace).Update(h.ctx, daemonset, h.Options.UpdateOptions)
}

// UpdateFromBytes update daemonset from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*appsv1.DaemonSet, error) {
	dsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	daemonset := &appsv1.DaemonSet{}
	err = json.Unmarshal(dsJson, daemonset)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(daemonset.Namespace) != 0 {
		namespace = daemonset.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.AppsV1().DaemonSets(namespace).Update(h.ctx, daemonset, h.Options.UpdateOptions)
}

// UpdateFromFile update daemonset from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*appsv1.DaemonSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update daemonset from file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*appsv1.DaemonSet, error) {
	return h.UpdateFromFile(filename)
}
