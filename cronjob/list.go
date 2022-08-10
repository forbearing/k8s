package cronjob

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// List list all cronjobs in the k8s cluster, it simply call `ListAll`.
func (h *Handler) List() ([]*batchv1.CronJob, error) {
	return h.ListAll()
}

// ListByLabel list cronjobs by labels.
// Multiple labels separated by comma(",") eg: "name=myapp,role=devops",
// and there is an "And" relationship between multiple labels.
func (h *Handler) ListByLabel(labels string) ([]*batchv1.CronJob, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	cjList, err := h.clientset.BatchV1().CronJobs(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(cjList), nil
}

// ListByField list cronjobs by field, work like `kubectl get xxx --field-selector=xxx`.
func (h *Handler) ListByField(field string) ([]*batchv1.CronJob, error) {
	fieldSelector, err := fields.ParseSelector(field)
	if err != nil {
		return nil, err
	}
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.FieldSelector = fieldSelector.String()

	cjList, err := h.clientset.BatchV1().CronJobs(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}
	return extractList(cjList), nil
}

// ListByNamespace list all cronjobs in the specified namespace.
func (h *Handler) ListByNamespace(namespace string) ([]*batchv1.CronJob, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all cronjobs in the k8s cluster.
func (h *Handler) ListAll() ([]*batchv1.CronJob, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}

// extractList
func extractList(cjList *batchv1.CronJobList) []*batchv1.CronJob {
	var objList []*batchv1.CronJob
	for i := range cjList.Items {
		objList = append(objList, &cjList.Items[i])
	}
	return objList
}
