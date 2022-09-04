package dynamic

import (
	"fmt"

	"github.com/forbearing/k8s/types"
	utilrestmapper "github.com/forbearing/k8s/util/restmapper"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

	gvr, isNamespaced, err := h.getGVRAndNamespaceScope()
	if err != nil {
		return nil, err
	}
	if isNamespaced {
		unstructList, err := h.dynamicClient.Resource(gvr).Namespace(h.namespace).List(h.ctx, *listOptions)
		if err != nil {
			return nil, err
		}
		return extractList(unstructList), nil
	}
	unstructList, err := h.dynamicClient.Resource(gvr).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(unstructList), nil
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

	gvr, isNamespaced, err := h.getGVRAndNamespaceScope()
	if err != nil {
		return nil, err
	}
	if isNamespaced {
		unstructList, err := h.dynamicClient.Resource(gvr).Namespace(h.namespace).List(h.ctx, *listOptions)
		if err != nil {
			return nil, err
		}
		return extractList(unstructList), nil
	}
	unstructList, err := h.dynamicClient.Resource(gvr).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(unstructList), nil
}

// ListByNamespace list all k8s objects in the specified namespace.
// It will return empty slice and error if this k8s object is cluster scope.
// Calling this method requires WithGVK() to explicitly specify GVK.
func (h *Handler) ListByNamespace(namespace string) ([]*unstructured.Unstructured, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""

	gvr, isNamespaced, err := h.getGVRAndNamespaceScope()
	if err != nil {
		return nil, err
	}
	if isNamespaced {
		unstructList, err := h.dynamicClient.Resource(gvr).Namespace(namespace).List(h.ctx, *listOptions)
		if err != nil {
			return nil, err
		}
		return extractList(unstructList), nil
	}
	return nil, fmt.Errorf("%s is not namespace-scoped k8s resource", gvr)
}

// ListAll list all k8s objects in the k8s cluster.
// Calling this method requires WithGVK() to explicitly specify GVK.
func (h *Handler) ListAll() ([]*unstructured.Unstructured, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""

	gvr, isNamespaced, err := h.getGVRAndNamespaceScope()
	if err != nil {
		return nil, err
	}
	if isNamespaced {
		unstructList, err := h.dynamicClient.Resource(gvr).Namespace(metav1.NamespaceAll).List(h.ctx, *listOptions)
		if err != nil {
			return nil, err
		}
		return extractList(unstructList), nil
	}
	unstructList, err := h.dynamicClient.Resource(gvr).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(unstructList), nil
}

// extractList
func extractList(unstructList *unstructured.UnstructuredList) []*unstructured.Unstructured {
	var objList []*unstructured.Unstructured
	for i := range unstructList.Items {
		objList = append(objList, &unstructList.Items[i])
	}
	return objList
}

func (h *Handler) getGVRAndNamespaceScope() (schema.GroupVersionResource, bool, error) {
	var (
		err          error
		gvr          schema.GroupVersionResource
		isNamespaced bool
	)

	if gvr, err = utilrestmapper.GVKToGVR(h.restMapper, h.gvk); err != nil {
		return gvr, isNamespaced, err
	}
	if isNamespaced, err = utilrestmapper.IsNamespaced(h.restMapper, h.gvk); err != nil {
		return gvr, isNamespaced, err
	}
	if h.gvk.Kind == types.KindJob || h.gvk.Kind == types.KindCronJob {
		h.SetPropagationPolicy("background")
	}

	return gvr, isNamespaced, nil
}
