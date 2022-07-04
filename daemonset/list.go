package daemonset

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListByLabel list daemonsets by labels.
func (h *Handler) ListByLabel(labels string) (*appsv1.DaemonSetList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.AppsV1().DaemonSets(h.namespace).List(h.ctx, *listOptions)
}

// List list daemonsets by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*appsv1.DaemonSetList, error) {
	return h.ListByLabel(labels)
}

// ListByNamespace list daemonsets by namespace.
func (h *Handler) ListByNamespace(namespace string) (*appsv1.DaemonSetList, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all daemonsets  in the k8s cluster.
func (h *Handler) ListAll() (*appsv1.DaemonSetList, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}
