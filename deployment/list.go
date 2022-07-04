package deployment

import (
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListByLabel list deployments by labels.
func (h *Handler) ListByLabel(labels string) (*appsv1.DeploymentList, error) {
	//h.Options.ListOptions.LabelSelector = labelSelector
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.AppsV1().Deployments(h.namespace).List(h.ctx, *listOptions)
}

// List list deployments by labels, alias to "ListByLabel".
func (h *Handler) List(label string) (*appsv1.DeploymentList, error) {
	return h.ListByLabel(label)
}

//// ListByNode list deployments by k8s node name
//// deployment not support list by k8s node name
//func (h *Handler) ListByNode(name string) (*appsv1.DeploymentList, error) {
//    // ParseSelector takes a string representing a selector and returns an
//    // object suitable for matching, or an error.
//    fieldSelector, err := fields.ParseSelector(fmt.Sprintf("spec.nodeName=%s", name))
//    if err != nil {
//        return nil, err
//    }
//    listOptions := h.Options.ListOptions.DeepCopy()
//    listOptions.FieldSelector = fieldSelector.String()

//    return h.clientset.AppsV1().Deployments(metav1.NamespaceAll).List(h.ctx, *listOptions)
//}

// ListByNamespace list all deployments in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) (*appsv1.DeploymentList, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all deployments in the k8s cluster.
func (h *Handler) ListAll() (*appsv1.DeploymentList, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}
