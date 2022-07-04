package pod

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// ListByLabel list pods by labels.
func (h *Handler) ListByLabel(labels string) (*corev1.PodList, error) {
	// TODO: 合并 ListOptions
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.CoreV1().Pods(h.namespace).List(h.ctx, *listOptions)
}

// List list pods by labels, alias to "ListByLabel"
func (h *Handler) List(labels string) (*corev1.PodList, error) {
	return h.ListByLabel(labels)
}

// ListByNode list all pods in k8s node where the pod is running.
func (h *Handler) ListByNode(name string) (*corev1.PodList, error) {
	// ParseSelector takes a string representing a selector and returns an
	// object suitable for matching, or an error.
	fieldSelector, err := fields.ParseSelector(fmt.Sprintf("spec.nodeName=%s", name))
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	return h.clientset.CoreV1().Pods(metav1.NamespaceAll).List(h.ctx, *listOptions)
}

// ListByNamespace list all pods in the specified namespace
func (h *Handler) ListByNamespace(namespace string) (*corev1.PodList, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all pods in k8s cluster where the pod is running.
func (h *Handler) ListAll() (*corev1.PodList, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}
