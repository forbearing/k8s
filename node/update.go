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

// Update updates node from type string, []byte, *corev1.Node,
// corev1.Node, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*corev1.Node, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *corev1.Node:
		return h.UpdateFromObject(val)
	case corev1.Node:
		return h.UpdateFromObject(&val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.UpdateFromObject(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates node from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.Node, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates node from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.Node, error) {
	nodeJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	node := &corev1.Node{}
	if err = json.Unmarshal(nodeJson, node); err != nil {
		return nil, err
	}
	return h.updateNode(node)
}

// UpdateFromObject updates node from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*corev1.Node, error) {
	node, ok := obj.(*corev1.Node)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Node")
	}
	return h.updateNode(node)
}

// UpdateFromUnstructured updates node from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*corev1.Node, error) {
	node := &corev1.Node{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), node)
	if err != nil {
		return nil, err
	}
	return h.updateNode(node)
}

// UpdateFromMap updates node from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*corev1.Node, error) {
	node := &corev1.Node{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, node)
	if err != nil {
		return nil, err
	}
	return h.updateNode(node)
}

// updateNode
func (h *Handler) updateNode(node *corev1.Node) (*corev1.Node, error) {
	node.ResourceVersion = ""
	node.UID = ""
	return h.clientset.CoreV1().Nodes().Update(h.ctx, node, h.Options.UpdateOptions)
}
