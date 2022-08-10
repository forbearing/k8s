package replicationcontroller

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all replicationcontrollers in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*corev1.ReplicationController, error) {
	return h.ListAll()
}

// ListByLabel list replicationcontrollers by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*corev1.ReplicationController, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	rcList, err := h.clientset.CoreV1().ReplicationControllers(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(rcList), nil
}

// ListByField list replicationcontrollers by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*corev1.ReplicationController, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	rcList, err := h.clientset.CoreV1().ReplicationControllers(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(rcList), nil
}

// ListByNamespace list all replicationcontrollers in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*corev1.ReplicationController, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all replicationcontrollers in the k8s cluster.
func (h *Handler) ListAll() ([]*corev1.ReplicationController, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(rcList *corev1.ReplicationControllerList) []*corev1.ReplicationController {
	var objList []*corev1.ReplicationController
	for i := range rcList.Items {
		objList = append(objList, &rcList.Items[i])
	}
	return objList
}
