package role

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all roles in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*rbacv1.Role, error) {
	return h.ListAll()
}

// ListByLabel list roles by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*rbacv1.Role, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	roleList, err := h.clientset.RbacV1().Roles(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(roleList), nil
}

// ListByField list roles by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*rbacv1.Role, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	roleList, err := h.clientset.RbacV1().Roles(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(roleList), nil
}

// ListByNamespace list all roles in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*rbacv1.Role, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all roles in the k8s cluster.
func (h *Handler) ListAll() ([]*rbacv1.Role, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(roleList *rbacv1.RoleList) []*rbacv1.Role {
	var objList []*rbacv1.Role
	for i := range roleList.Items {
		objList = append(objList, &roleList.Items[i])
	}
	return objList
}
