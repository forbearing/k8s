package clusterrolebinding

import (
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all clusterrolebindings in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*rbacv1.ClusterRoleBinding, error) {
	return h.ListAll()
}

// ListByLabel list clusterrolebindings by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*rbacv1.ClusterRoleBinding, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	crbList, err := h.clientset.RbacV1().ClusterRoleBindings().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(crbList), nil
}

// ListByField list clusterrolebindings by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*rbacv1.ClusterRoleBinding, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	crbList, err := h.clientset.RbacV1().ClusterRoleBindings().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(crbList), nil
}

// ListAll list all clusterrolebindings in the k8s cluster.
func (h *Handler) ListAll() ([]*rbacv1.ClusterRoleBinding, error) {
	return h.ListByLabel("")
}

// extractList
func extractList(crbList *rbacv1.ClusterRoleBindingList) []*rbacv1.ClusterRoleBinding {
	var objList []*rbacv1.ClusterRoleBinding
	for i := range crbList.Items {
		objList = append(objList, &crbList.Items[i])
	}
	return objList
}
