package role

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListByLabel list roles by labels.
func (h *Handler) ListByLabel(labels string) (*rbacv1.RoleList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.RbacV1().Roles(h.namespace).List(h.ctx, *listOptions)
}

// List list roles by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*rbacv1.RoleList, error) {
	return h.ListByLabel(labels)
}

// ListByNamespace list roles by namespace.
func (h *Handler) ListByNamespace(namespace string) (*rbacv1.RoleList, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all roles in the k8s cluster.
func (h *Handler) ListAll() (*rbacv1.RoleList, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}
