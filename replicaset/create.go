package replicaset

import (
	"encoding/json"
	"io/ioutil"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create replicaset from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*appsv1.ReplicaSet, error) {
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

	return h.clientset.AppsV1().ReplicaSets(namespace).Create(h.ctx, rs, h.Options.CreateOptions)
}

// CreateFromBytes create replicaset from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*appsv1.ReplicaSet, error) {
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

	return h.clientset.AppsV1().ReplicaSets(namespace).Create(h.ctx, replicaset, h.Options.CreateOptions)
}

// CreateFromFile create replicaset from yaml file.
func (h *Handler) CreateFromFile(filename string) (*appsv1.ReplicaSet, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create replicaset from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*appsv1.ReplicaSet, error) {
	return h.CreateFromFile(filename)
}
