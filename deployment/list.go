package deployment

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all deployments in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*appsv1.Deployment, error) {
	return h.ListAll()
}

// ListByLabel list deployments by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*appsv1.Deployment, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	deployList, err := h.clientset.AppsV1().Deployments(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(deployList), nil
}

// ListByField list deployments by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*appsv1.Deployment, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	deployList, err := h.clientset.AppsV1().Deployments(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(deployList), nil
}

// ListByNamespace list all deployments in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*appsv1.Deployment, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all deployments in the k8s cluster.
func (h *Handler) ListAll() ([]*appsv1.Deployment, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(deployList *appsv1.DeploymentList) []*appsv1.Deployment {
	var objList []*appsv1.Deployment
	for i := range deployList.Items {
		objList = append(objList, &deployList.Items[i])
	}
	return objList
}
