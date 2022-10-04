package dynamic

import (
	"fmt"

	"github.com/forbearing/k8s/types"
	utilrestmapper "github.com/forbearing/k8s/util/restmapper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all k8s objects in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*unstructured.Unstructured, error) {
	return h.ListAll()
}

// ListByLabel list k8s objects by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
// Calling this method requires WithGVK() to explicitly specify GVK.
func (h *Handler) ListByLabel(labels string) ([]*unstructured.Unstructured, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels

	if err := h.getGVRAndNamespaceScope(); err != nil {
		return nil, err
	}
	if h.isNamespaced {
		return extractList(h.dynamicClient.Resource(h.gvr).Namespace(h.namespace).List(h.ctx, *listOptions))
	}
	return extractList(h.dynamicClient.Resource(h.gvr).List(h.ctx, *listOptions))
}

// ListByField list k8s objects by field, work like `kubectl get xxx --field-selector=xxx`.
// Calling this method requires WithGVK() to explicitly specify GVK.
func (h *Handler) ListByField(field string) ([]*unstructured.Unstructured, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	if err := h.getGVRAndNamespaceScope(); err != nil {
		return nil, err
	}
	if h.isNamespaced {
		return extractList(h.dynamicClient.Resource(h.gvr).Namespace(h.namespace).List(h.ctx, *listOptions))
	}
	return extractList(h.dynamicClient.Resource(h.gvr).List(h.ctx, *listOptions))
}

// ListByNamespace list all k8s objects in the specified namespace.
// It will return empty slice and error if this k8s object is cluster scope.
// Calling this method requires WithGVK() to explicitly specify GVK.
func (h *Handler) ListByNamespace(namespace string) ([]*unstructured.Unstructured, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""

	if err := h.getGVRAndNamespaceScope(); err != nil {
		return nil, err
	}
	if h.isNamespaced {
		return extractList(h.dynamicClient.Resource(h.gvr).Namespace(namespace).List(h.ctx, *listOptions))
	}
	return nil, fmt.Errorf("%s is not namespace-scoped k8s resource", h.gvr)
}

// ListAll list all k8s objects in the k8s cluster.
// Calling this method requires WithGVK() to explicitly specify GVK.
func (h *Handler) ListAll() ([]*unstructured.Unstructured, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""

	if err := h.getGVRAndNamespaceScope(); err != nil {
		return nil, err
	}
	if h.isNamespaced {
		return extractList(h.dynamicClient.Resource(h.gvr).Namespace(metav1.NamespaceAll).List(h.ctx, *listOptions))
	}
	return extractList(h.dynamicClient.Resource(h.gvr).List(h.ctx, *listOptions))
}

// extractList
func extractList(unstructList *unstructured.UnstructuredList, err error) ([]*unstructured.Unstructured, error) {
	if err != nil {
		return nil, err
	}
	var objList []*unstructured.Unstructured
	for i := range unstructList.Items {
		objList = append(objList, &unstructList.Items[i])
	}
	return objList, nil
}

func (h *Handler) getGVRAndNamespaceScope() error {
	var err error
	if h.gvr, err = utilrestmapper.GVKToGVR(h.restMapper, h.gvk); err != nil {
		return err
	}
	if h.isNamespaced, err = utilrestmapper.IsNamespaced(h.restMapper, h.gvk); err != nil {
		return err
	}
	if h.gvk.Kind == types.KindJob || h.gvk.Kind == types.KindCronJob {
		h.SetPropagationPolicy("background")
	}

	return nil
}
