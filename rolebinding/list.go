package rolebinding

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all rolebindings in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*rbacv1.RoleBinding, error) {
	return h.ListAll()
}

// ListByLabel list rolebindings by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*rbacv1.RoleBinding, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	rbList, err := h.clientset.RbacV1().RoleBindings(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(rbList), nil
}

// ListByField list rolebindings by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*rbacv1.RoleBinding, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	rbList, err := h.clientset.RbacV1().RoleBindings(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(rbList), nil
}

// ListByNamespace list all rolebindings in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*rbacv1.RoleBinding, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all rolebindings in the k8s cluster.
func (h *Handler) ListAll() ([]*rbacv1.RoleBinding, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(rbList *rbacv1.RoleBindingList) []*rbacv1.RoleBinding {
	var objList []*rbacv1.RoleBinding
	for i := range rbList.Items {
		objList = append(objList, &rbList.Items[i])
	}
	return objList
}
