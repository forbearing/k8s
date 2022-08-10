package networkpolicy

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all networkpolicies in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*networkingv1.NetworkPolicy, error) {
	return h.ListAll()
}

// ListByLabel list networkpolicies by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*networkingv1.NetworkPolicy, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	netpolList, err := h.clientset.NetworkingV1().NetworkPolicies(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(netpolList), nil
}

// ListByField list networkpolicies by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*networkingv1.NetworkPolicy, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	netpolList, err := h.clientset.NetworkingV1().NetworkPolicies(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(netpolList), nil
}

// ListByNamespace list all networkpolicies in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*networkingv1.NetworkPolicy, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all networkpolicies in the k8s cluster.
func (h *Handler) ListAll() ([]*networkingv1.NetworkPolicy, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(netpolList *networkingv1.NetworkPolicyList) []*networkingv1.NetworkPolicy {
	var objList []*networkingv1.NetworkPolicy
	for i := range netpolList.Items {
		objList = append(objList, &netpolList.Items[i])
	}
	return objList
}
