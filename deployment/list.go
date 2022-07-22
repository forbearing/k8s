package deployment

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// List list deployments by labels, alias to "ListByLabel".
func (h *Handler) List(label string) (*appsv1.DeploymentList, error) {
	return h.ListByLabel(label)
}

// ListByLabel list deployments by labels.
func (h *Handler) ListByLabel(labels string) (*appsv1.DeploymentList, error) {
	//h.Options.ListOptions.LabelSelector = labelSelector
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.AppsV1().Deployments(h.namespace).List(h.ctx, *listOptions)
}

// ListByNamespace list all deployments in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) (*appsv1.DeploymentList, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all deployments in the k8s cluster.
func (h *Handler) ListAll() (*appsv1.DeploymentList, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}
