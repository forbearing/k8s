package dynamic

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

func (h *Handler) List(label string) (*unstructured.UnstructuredList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = label
	if h.namespacedResource {
		h.dynamicClient.Resource(h.gvr()).Namespace(h.namespace).List(h.ctx, *listOptions)
	}
	return h.dynamicClient.Resource(h.gvr()).Namespace(h.namespace).List(h.ctx, *listOptions)
}
