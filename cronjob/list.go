package cronjob

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListByLabel list cronjobs by labels.
func (h *Handler) ListByLabel(labels string) (*batchv1.CronJobList, error) {
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = labels
	return h.clientset.BatchV1().CronJobs(h.namespace).List(h.ctx, *listOptions)
}

// List list cronjobs by labels, alias to ListByLabel.
func (h *Handler) List(labels string) (*batchv1.CronJobList, error) {
	return h.ListByLabel(labels)
}

// ListByNamespace list cronjobs by namespace.
func (h *Handler) ListByNamespace(namespace string) (*batchv1.CronJobList, error) {
	return h.WithNamespace(namespace).ListByLabel("")
}

// ListAll list all cronjobs in the k8s cluster.
func (h *Handler) ListAll(namespace string) (*batchv1.CronJobList, error) {
	return h.WithNamespace(metav1.NamespaceAll).ListByLabel("")
}
