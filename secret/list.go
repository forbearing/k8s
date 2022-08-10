package secret

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all cecrets in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*corev1.Secret, error) {
	return h.ListAll()
}

// ListByLabel list cecrets by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*corev1.Secret, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	secretList, err := h.clientset.CoreV1().Secrets(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(secretList), nil
}

// ListByField list cecrets by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*corev1.Secret, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	secretList, err := h.clientset.CoreV1().Secrets(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(secretList), nil
}

// ListByNamespace list all cecrets in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*corev1.Secret, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all cecrets in the k8s cluster.
func (h *Handler) ListAll() ([]*corev1.Secret, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(secretList *corev1.SecretList) []*corev1.Secret {
	var objList []*corev1.Secret
	for i := range secretList.Items {
		objList = append(objList, &secretList.Items[i])
	}
	return objList
}
