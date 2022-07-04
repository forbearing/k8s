package replicaset

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update replicaset from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*appsv1.ReplicaSet, error) {
	rs := &appsv1.ReplicaSet{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, rs)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(rs.Namespace) != 0 {
		namespace = rs.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.AppsV1().ReplicaSets(namespace).Update(h.ctx, rs, h.Options.UpdateOptions)
}

// UpdateFromBytes update replicaset from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*appsv1.ReplicaSet, error) {
	dsJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	replicaset := &appsv1.ReplicaSet{}
	err = json.Unmarshal(dsJson, replicaset)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(replicaset.Namespace) != 0 {
		namespace = replicaset.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.AppsV1().ReplicaSets(namespace).Update(h.ctx, replicaset, h.Options.UpdateOptions)
}

// UpdateFromFile update replicaset from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*appsv1.ReplicaSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update replicaset from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*appsv1.ReplicaSet, error) {
	return h.UpdateFromFile(filename)
}
