package node

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates node from type string, []byte, *corev1.Node,
// corev1.Node, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*corev1.Node, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *corev1.Node:
		return h.CreateFromObject(val)
	case corev1.Node:
		return h.CreateFromObject(&val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.CreateFromObject(val)
	default:
		return nil, ErrInvalidCreateType
	}
}

// CreateFromFile creates node from yaml or json file.
func (h *Handler) CreateFromFile(filename string) (*corev1.Node, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates node from bytes data.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.Node, error) {
	nodeJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	node := &corev1.Node{}
	if err = json.Unmarshal(nodeJson, node); err != nil {
		return nil, err
	}
	return h.createNode(node)
}

// CreateFromObject creates node from metav1.Object or runtime.Object.
func (h *Handler) CreateFromObject(obj interface{}) (*corev1.Node, error) {
	node, ok := obj.(*corev1.Node)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Node")
	}
	return h.createNode(node)
}

// CreateFromUnstructured creates node from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*corev1.Node, error) {
	node := &corev1.Node{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), node)
	if err != nil {
		return nil, err
	}
	return h.createNode(node)
}

// CreateFromMap creates node from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*corev1.Node, error) {
	node := &corev1.Node{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, node)
	if err != nil {
		return nil, err
	}
	return h.createNode(node)
}

// createNode
func (h *Handler) createNode(node *corev1.Node) (*corev1.Node, error) {
	node.ResourceVersion = ""
	node.UID = ""
	return h.clientset.CoreV1().Nodes().Create(h.ctx, node, h.Options.CreateOptions)
}
