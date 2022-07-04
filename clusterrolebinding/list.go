package clusterrolebinding

import rbacv1 "k8s.io/api/rbac/v1"

// ListByLabel list clusterrolebindings by labels.
func (h *Handler) ListByLabel(labels string) (*rbacv1.ClusterRoleBindingList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.RbacV1().ClusterRoleBindings().List(h.ctx, *listOptions)
}

// List list clusterrolebindings by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*rbacv1.ClusterRoleBindingList, error) {
	return h.ListByLabel(labels)
}

// ListAll list all clusterrolebindings in the k8s cluster.
func (h *Handler) ListAll() (*rbacv1.ClusterRoleBindingList, error) {
	return h.ListByLabel("")
}
