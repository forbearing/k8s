package statefulset

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create statefulset from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*appsv1.StatefulSet, error) {
	sts := &appsv1.StatefulSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, sts)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(sts.Namespace) != 0 {
		namespace = sts.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.AppsV1().StatefulSets(namespace).Create(h.ctx, sts, h.Options.CreateOptions)
}

// CreateFromBytes create statefulset from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*appsv1.StatefulSet, error) {
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

	return h.clientset.AppsV1().StatefulSets(namespace).Create(h.ctx, sts, h.Options.CreateOptions)
}

// CreateFromFile create statefulset from yaml file.
func (h *Handler) CreateFromFile(filename string) (*appsv1.StatefulSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create statefulset from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*appsv1.StatefulSet, error) {
	return h.CreateFromFile(filename)
}
