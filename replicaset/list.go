package replicaset

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListByLabel list replicasets by labels.
func (h *Handler) ListByLabel(labels string) (*appsv1.ReplicaSetList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.AppsV1().ReplicaSets(h.namespace).List(h.ctx, *listOptions)
}

// List list replicasets by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*appsv1.ReplicaSetList, error) {
	return h.ListByLabel(labels)
}

// ListByNamespace list replicasets by namespace.
func (h *Handler) ListByNamespace(namespace string) (*appsv1.ReplicaSetList, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all replicasets in the k8s cluster.
func (h *Handler) ListAll() (*appsv1.ReplicaSetList, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}
