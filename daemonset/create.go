package daemonset

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create daemonset from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*appsv1.DaemonSet, error) {
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

	return h.clientset.AppsV1().DaemonSets(namespace).Create(h.ctx, daemonset, h.Options.CreateOptions)
}

// CreateFromBytes create daemonset from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*appsv1.DaemonSet, error) {
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

	return h.clientset.AppsV1().DaemonSets(namespace).Create(h.ctx, daemonset, h.Options.CreateOptions)
}

// CreateFromFile create daemonset from yaml file.
func (h *Handler) CreateFromFile(filename string) (*appsv1.DaemonSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create daemonset from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*appsv1.DaemonSet, error) {
	return h.CreateFromFile(filename)
}
