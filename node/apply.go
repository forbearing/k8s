package node

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies node from type string, []byte, *corev1.Node,
// corev1.Node, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*corev1.Node, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *corev1.Node:
		return h.ApplyFromObject(val)
	case corev1.Node:
		return h.ApplyFromObject(&val)
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	case metav1.Object, runtime.Object:
		return h.ApplyFromObject(val)
	default:
		return nil, ErrInvalidApplyType
	}
}

// ApplyFromFile applies node from yaml or json file.
func (h *Handler) ApplyFromFile(filename string) (node *corev1.Node, err error) {
	node, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if node already exist, update it.
		node, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply node from bytes data.
func (h *Handler) ApplyFromBytes(data []byte) (node *corev1.Node, err error) {
	node, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		node, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies node from metav1.Object or runtime.Object.
func (h *Handler) ApplyFromObject(obj interface{}) (*corev1.Node, error) {
	node, ok := obj.(*corev1.Node)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Node")
	}
	return h.applyNode(node)
}

// ApplyFromUnstructured applies node from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*corev1.Node, error) {
	node := &corev1.Node{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), node)
	if err != nil {
		return nil, err
	}
	return h.applyNode(node)
}

// ApplyFromMap applies node from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*corev1.Node, error) {
	node := &corev1.Node{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, node)
	if err != nil {
		return nil, err
	}
	return h.applyNode(node)
}

// applyNode
func (h *Handler) applyNode(node *corev1.Node) (*corev1.Node, error) {
	_, err := h.createNode(node)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateNode(node)
	}
	return node, err
}
