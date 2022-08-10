package clusterrole

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all clusterroles in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*rbacv1.ClusterRole, error) {
	return h.ListAll()
}

// ListByLabel list clusterroles by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*rbacv1.ClusterRole, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	crList, err := h.clientset.RbacV1().ClusterRoles().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(crList), nil
}

// ListByField list clusterroles by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*rbacv1.ClusterRole, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	crList, err := h.clientset.RbacV1().ClusterRoles().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(crList), nil
}

// ListAll list all clusterroles in the k8s cluster.
func (h *Handler) ListAll() ([]*rbacv1.ClusterRole, error) {
	return h.ListByLabel("")
}

// extractList
func extractList(crList *rbacv1.ClusterRoleList) []*rbacv1.ClusterRole {
	var objList []*rbacv1.ClusterRole
	for i := range crList.Items {
		objList = append(objList, &crList.Items[i])
	}
	return objList
}
