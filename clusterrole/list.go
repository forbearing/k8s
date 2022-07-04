package clusterrole

import rbacv1 "k8s.io/api/rbac/v1"

// ListByLabel list clusterroles by labels.
func (h *Handler) ListByLabel(labels string) (*rbacv1.ClusterRoleList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.RbacV1().ClusterRoles().List(h.ctx, *listOptions)
}

// List list clusterroles by labels, alias to "ListByLabel".
func (h *Handler) List(labels string) (*rbacv1.ClusterRoleList, error) {
	return h.ListByLabel(labels)
}

// ListAll list all clusterroles in the k8s cluster.
func (h *Handler) ListAll(labels string) (*rbacv1.ClusterRoleList, error) {
	return h.ListByLabel("")
}
