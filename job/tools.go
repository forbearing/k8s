package job

import (
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

var ERR_TYPE = fmt.Errorf("type must be *batchv1.Job, batchv1.Job or string")

//type JobController struct {
//    Labels             map[string]string `json:"labels"`
//    CreationTimestamp  metav1.Time       `json:"creationTimestamp"`
//    LastScheduleTime   metav1.Time       `json:"lastScheduleTime"`
//    LastSuccessfulTime metav1.Time       `json:"lastSuccessfulTime"`

//    metav1.OwnerReference `json:"ownerReference"`
//}
type JobController struct {
	metav1.OwnerReference `json:"ownerReference"`
}

// IsCompleted will check if the job was successfully scheduled and run to completed
// job 成功调度并且其 pod 成功执行
func (h *Handler) IsCompleted(name string) bool {
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

// IsFailed will check if the job was successfully scheduled but run to failed.
func (h *Handler) IsFailed(name string) bool {
	job, err := h.Get(name)
	if err != nil {
		return false
	}

	for _, cond := range job.Status.Conditions {
		if cond.Status == corev1.ConditionTrue && cond.Type == batchv1.JobFailed {
			return true
		}
	}
	return false
}

// IsSuspended will check if the job was successfully scheduled but the job was suspended.
func (h *Handler) IsSuspended(name string) bool {
	job, err := h.Get(name)
	if err != nil {
		return false
	}

	for _, cond := range job.Status.Conditions {
		if cond.Status == corev1.ConditionTrue && cond.Type == batchv1.JobSuspended {
			return true
		}
	}
	return false
}

// IsFinished will check the job was successfully scheduled, it doesn't matter
// if the job runs to completion or fails.
// job 成功调度, pod 不再产生, 比如 job 设置了 backoffLimit: 3, 则限制 job 产生的
// pod 最多失败3次, 超过3次, job 不再创建 pod 来做执行任务. 这样也是 Finished
// 总结就是: job 不再创建 pod 来执行任务, 不管任务是否成功或失败
// 如果执行任务成功,  job 只会创建一次 pod, 如果任务执行失败,创建多少次 pod 取决于
// backoffLimit 的设置.
func (h *Handler) IsFinished(name string) bool {
	// 1. job not exist, return true
	job, err := h.Get(name)
	if err != nil {
		return true
	}
	// 1. if job complete and status is true, return true.
	// 2. if job failed and status is true, return true.
	for _, cond := range job.Status.Conditions {
		//if cond.Status == corev1.ConditionTrue && cond.Type == batchv1.JobComplete {
		//    return true
		//}
		//if cond.Status == corev1.ConditionTrue && cond.Type == batchv1.JobFailed {
		//    return true
		//}
		if (cond.Type == batchv1.JobComplete || cond.Type == batchv1.JobFailed) && cond.Status == corev1.ConditionTrue {
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
	if h.IsFinished(name) {
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
				if h.IsFinished(name) {
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

func (h *Handler) getController(jobObj *batchv1.Job) (*JobController, error) {
	ownerRef := metav1.GetControllerOf(jobObj)
	if ownerRef == nil {
		return nil, fmt.Errorf("the job %q doesn't have controller", jobObj.Name)
	}
	oc := JobController{OwnerReference: *ownerRef}
	return &oc, nil
}

// GetController returns a JobController object by job name if the controllee(job) has a controller.
func (h *Handler) GetController(object interface{}) (*JobController, error) {
	switch val := object.(type) {
	case string:
		j, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getController(j)
	case *batchv1.Job:
		return h.getController(val)
	case batchv1.Job:
		return h.getController(&val)
	default:
		return nil, ERR_TYPE
	}
}

// DurationOfStarted returns the duration from job start time to now.
func (h *Handler) DurationOfStarted(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		j, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		if j.Status.StartTime == nil {
			return time.Duration(int64(0)), fmt.Errorf("job started time not found")
		}
		return time.Now().Sub(j.Status.StartTime.Time), nil
	case *batchv1.Job:
		if val.Status.StartTime == nil {
			return time.Duration(int64(0)), fmt.Errorf("job started time not found")
		}
		return time.Now().Sub(val.Status.StartTime.Time), nil
	case batchv1.Job:
		if val.Status.StartTime == nil {
			return time.Duration(int64(0)), fmt.Errorf("job started time not found")
		}
		return time.Now().Sub(val.Status.StartTime.Time), nil
	default:
		return time.Duration(int64(0)), ERR_TYPE
	}
}

// DurationOfCompletion returns the duration from job start time to now.
func (h *Handler) DurationOfCompleted(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		j, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		if j.Status.CompletionTime == nil {
			return time.Duration(int64(0)), fmt.Errorf("job completed time not found")
		}
		return time.Now().Sub(j.Status.CompletionTime.Time), nil
	case *batchv1.Job:
		if val.Status.CompletionTime == nil {
			return time.Duration(int64(0)), fmt.Errorf("job completed time not found")
		}
		return time.Now().Sub(val.Status.CompletionTime.Time), nil
	case batchv1.Job:
		if val.Status.CompletionTime == nil {
			return time.Duration(int64(0)), fmt.Errorf("job completed time not found")
		}
		return time.Now().Sub(val.Status.CompletionTime.Time), nil
	default:
		return time.Duration(int64(0)), ERR_TYPE
	}
}

// NumActive returns the number of pending or running pod created by given job.
func (h *Handler) NumActive(object interface{}) (int32, error) {
	switch val := object.(type) {
	case string:
		j, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return j.Status.Active, nil
	case *batchv1.Job:
		return val.Status.Active, nil
	case batchv1.Job:
		return val.Status.Active, nil
	default:
		return 0, nil
	}
}

// NumSucceeded returns the number of pods created by given job which reach phase Succeeded.
func (h *Handler) NumSucceeded(object interface{}) (int32, error) {
	switch val := object.(type) {
	case string:
		j, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return j.Status.Succeeded, nil
	case *batchv1.Job:
		return val.Status.Succeeded, nil
	case batchv1.Job:
		return val.Status.Succeeded, nil
	default:
		return 0, nil
	}
}

// NumFailed returns the number of pods created by given job which reach phase failed.
func (h *Handler) NumFailed(object interface{}) (int32, error) {
	switch val := object.(type) {
	case string:
		j, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return j.Status.Failed, nil
	case *batchv1.Job:
		return val.Status.Failed, nil
	case batchv1.Job:
		return val.Status.Failed, nil
	default:
		return 0, nil
	}
}

// GetAge get job age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		j, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		return time.Now().Sub(j.CreationTimestamp.Time), nil
	case *batchv1.Job:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case batchv1.Job:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(int64(0)), ERR_TYPE
	}
}
