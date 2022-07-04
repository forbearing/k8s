package rolebinding

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListByLabel list rolebindings by labels.
func (h *Handler) ListByLabel(labels string) (*rbacv1.RoleBindingList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.RbacV1().RoleBindings(h.namespace).List(h.ctx, *listOptions)
}

// List list rolebindings by labels, alias to  "ListByLabel".
func (h *Handler) List(labels string) (*rbacv1.RoleBindingList, error) {
	return h.ListByLabel(labels)
}

// ListByNamespace list rolebindings by namespace.
func (h *Handler) ListByNamespace(namespace string) (*rbacv1.RoleBindingList, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all rolebindings in the k8s cluster.
func (h *Handler) ListAll() (*rbacv1.RoleBindingList, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}
