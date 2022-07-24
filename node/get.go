package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets node from type string, []byte, *corev1.Node,
// corev1.Node, runtime.Object or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a node from file path.
func (h *Handler) Get(obj interface{}) (*corev1.Node, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *corev1.Node:
		return h.GetFromObject(val)
	case corev1.Node:
		return h.GetFromObject(&val)
	case map[string]interface{}:
		return h.GetFromUnstructured(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets node by name.
func (h *Handler) GetByName(name string) (*corev1.Node, error) {
	return h.clientset.CoreV1().Nodes().Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets node from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.Node, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets node from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.Node, error) {
	nodeJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	node := &corev1.Node{}
	err = json.Unmarshal(nodeJson, node)
	if err != nil {
		return nil, err
	}
	return h.getNode(node)
}

// GetFromObject gets node from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*corev1.Node, error) {
	node, ok := obj.(*corev1.Node)
	if !ok {
		return nil, fmt.Errorf("object is not *corev1.Node")
	}
	return h.getNode(node)
}

// GetFromUnstructured gets node from map[string]interface{}.
func (h *Handler) GetFromUnstructured(u map[string]interface{}) (*corev1.Node, error) {
	node := &corev1.Node{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, node)
	if err != nil {
		return nil, err
	}
	return h.getNode(node)
}

// getNode
// It's necessary to get a new node resource from a old node resource,
// because old node usually don't have node.Status field.
func (h *Handler) getNode(node *corev1.Node) (*corev1.Node, error) {
	return h.clientset.CoreV1().Nodes().Get(h.ctx, node.Name, h.Options.GetOptions)
}
