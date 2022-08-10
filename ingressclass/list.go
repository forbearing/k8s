package ingressclass

import (
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all ingressclasses in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*networkingv1.IngressClass, error) {
	return h.ListAll()
}

// ListByLabel list ingressclasses by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*networkingv1.IngressClass, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	ingcList, err := h.clientset.NetworkingV1().IngressClasses().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(ingcList), nil
}

// ListByField list ingressclasses by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*networkingv1.IngressClass, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	ingcList, err := h.clientset.NetworkingV1().IngressClasses().List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(ingcList), nil
}

// ListAll list all ingressclasses in the k8s cluster.
func (h *Handler) ListAll() ([]*networkingv1.IngressClass, error) {
	return h.ListByLabel("")
}

// extractList
func extractList(ingcList *networkingv1.IngressClassList) []*networkingv1.IngressClass {
	var objList []*networkingv1.IngressClass
	for i := range ingcList.Items {
		objList = append(objList, &ingcList.Items[i])
	}
	return objList
}
