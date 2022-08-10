package node

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all nodes in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*corev1.Node, error) {
	return h.ListAll()
}

// ListByLabel list nodes by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*corev1.Node, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	nodeList, err := h.clientset.CoreV1().Nodes().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(nodeList), nil
}

// ListByField list nodes by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*corev1.Node, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	nodeList, err := h.clientset.CoreV1().Nodes().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(nodeList), nil
}

// ListAll list all nodes in the k8s cluster.
func (h *Handler) ListAll() ([]*corev1.Node, error) {
	return h.ListByLabel("")
}

// extractList
func extractList(nodeList *corev1.NodeList) []*corev1.Node {
	var objList []*corev1.Node
	for i := range nodeList.Items {
		objList = append(objList, &nodeList.Items[i])
	}
	return objList
}
