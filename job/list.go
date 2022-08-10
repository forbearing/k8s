package job

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all jobs in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*batchv1.Job, error) {
	return h.ListAll()
}

// ListByLabel list jobs by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*batchv1.Job, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	jobList, err := h.clientset.BatchV1().Jobs(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(jobList), nil
}

// ListByField list jobs by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*batchv1.Job, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	jobList, err := h.clientset.BatchV1().Jobs(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(jobList), nil
}

// ListByNamespace list all jobs in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*batchv1.Job, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all jobs in the k8s cluster.
func (h *Handler) ListAll() ([]*batchv1.Job, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(jobList *batchv1.JobList) []*batchv1.Job {
	var objList []*batchv1.Job
	for i := range jobList.Items {
		objList = append(objList, &jobList.Items[i])
	}
	return objList
}
