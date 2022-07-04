package statefulset

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update statefulset from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*appsv1.StatefulSet, error) {
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

	return h.clientset.AppsV1().StatefulSets(namespace).Update(h.ctx, sts, h.Options.UpdateOptions)
}

// UpdateFromBytes update statefulset from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*appsv1.StatefulSet, error) {
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

	return h.clientset.AppsV1().StatefulSets(namespace).Update(h.ctx, sts, h.Options.UpdateOptions)
}

// UpdateFromFile update statefulset from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*appsv1.StatefulSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update statefulset from file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*appsv1.StatefulSet, error) {
	return h.UpdateFromFile(filename)
}
