package job

import (
	"fmt"
	"time"

	"github.com/forbearing/k8s/cronjob"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// GetController returns a JobController object by job name if the controllee(job) has a controller.
func (h *Handler) GetController(name string) (*JobController, error) {
	if len(name) == 0 {
		return nil, fmt.Errorf("not set the job name")
	}
	job, err := h.Get(name)
	if err != nil {
		return nil, err
	}
	ownerRef := metav1.GetControllerOf(job)
	if ownerRef == nil {
		return nil, fmt.Errorf("the job %q doesn't have controller", name)
	}
	oc := JobController{OwnerReference: *ownerRef}

	// new a cronjob handler
	cronjobHandler, err := cronjob.New(h.ctx, h.namespace, h.kubeconfig)
	if err != nil {
		return nil, err
	}
	cronjob, err := cronjobHandler.Get(ownerRef.Name)
	if err != nil {
		return nil, err
	}

	oc.Labels = cronjob.Labels
	oc.CreationTimestamp = cronjob.ObjectMeta.CreationTimestamp
	oc.LastScheduleTime = *(cronjob.Status.LastScheduleTime)
	oc.LastSuccessfulTime = *(cronjob.Status.LastSuccessfulTime)
	return &oc, nil
}

// IsComplete check job if is completion
func (h *Handler) IsComplete(name string) bool {
	// if job not exist, return false
	job, err := h.Get(name)
	if err != nil {
		return false
	}

	for _, cond := range job.Status.Conditions {
		if cond.Status == corev1.ConditionTrue && cond.Type == batchv1.JobComplete {
			return true
		}
	}

	return false
}

// IsFinish check job if is condition is
// job finished means that the job condition is "complete" or "failed"
func (h *Handler) IsFinish(name string) bool {
	// 1. job not exist, return true
	job, err := h.Get(name)
	if err != nil {
		return true
	}
	// 2. if job complete return true
	// 3. if job failed return true
	// 4. all other job condition return false
	for _, cond := range job.Status.Conditions {
		if cond.Status == corev1.ConditionTrue && cond.Type == batchv1.JobComplete {
			return true
		}
		if cond.Status == corev1.ConditionTrue && cond.Type == batchv1.JobFailed {
			return true
		}
	}
	return false
}

// WaitFinish wait job status to be "true"
func (h *Handler) WaitFinish(name string) (err error) {
	var (
		watcher watch.Interface
		timeout = int64(0)
	)
	if h.IsFinish(name) {
		return
	}

	for {
		listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
		listOptions.TimeoutSeconds = &timeout
		watcher, err = h.clientset.BatchV1().Jobs(h.namespace).Watch(h.ctx, listOptions)
		if err != nil {
			return
		}
		for event := range watcher.ResultChan() {
			switch event.Type {
			case watch.Modified:
				if h.IsFinish(name) {
					watcher.Stop()
					return
				}
			case watch.Deleted:
				watcher.Stop()
				return fmt.Errorf("%s deleted", name)
			}
		}
	}
}

// WaitNotExist wait job not exist
func (h *Handler) WaitNotExist(name string) (err error) {
	var (
		watcher watch.Interface
		timeout = int64(0)
	)
	_, err = h.Get(name)
	if err != nil { // job not exist
		return nil
	}
	for {
		listOptions := metav1.SingleObject(metav1.ObjectMeta{Name: name, Namespace: h.namespace})
		listOptions.TimeoutSeconds = &timeout
		watcher, err = h.clientset.BatchV1().Jobs(h.namespace).Watch(h.ctx, listOptions)
		if err != nil {
			return
		}
		for event := range watcher.ResultChan() {
			switch event.Type {
			case watch.Deleted:
				for {
					if _, err := h.Get(name); err != nil { // job not exist
						break
					}
					time.Sleep(time.Millisecond * 500)
				}
				watcher.Stop()
				return
			}
		}
	}
}
