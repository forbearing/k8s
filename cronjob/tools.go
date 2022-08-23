package cronjob

import (
	"fmt"
	"time"

	batchv1 "k8s.io/api/batch/v1"
)

// GetJobs get all jobs created by the cronjob.
func (h *Handler) GetJobs(object interface{}) ([]batchv1.Job, error) {
	switch val := object.(type) {
	case string:
		cj, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getJobs(cj)
	case *batchv1.CronJob:
		return h.getJobs(val)
	case batchv1.CronJob:
		return h.getJobs(&val)
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getJobs(cj *batchv1.CronJob) ([]batchv1.Job, error) {
	// list all job in the same namespace as the cronjob
	listOptions := h.Options.ListOptions.DeepCopy()
	listOptions.LabelSelector = ""
	jobList, err := h.clientset.BatchV1().Jobs(h.namespace).List(h.ctx, *listOptions)
	if err != nil {
		return nil, err
	}

	var jl []batchv1.Job
	for _, j := range jobList.Items {
		for _, ownerRef := range j.OwnerReferences {
			if ownerRef.Name == cj.Name {
				jl = append(jl, j)
			}
		}
	}
	return jl, nil
}

// NumActive returns the number of running job created by cronjob.
func (h *Handler) NumActive(object interface{}) (int, error) {
	switch val := object.(type) {
	case string:
		cj, err := h.Get(val)
		if err != nil {
			return 0, err
		}
		return len(cj.Status.Active), nil
	case *batchv1.CronJob:
		return len(val.Status.Active), nil
	case batchv1.CronJob:
		return len(val.Status.Active), nil
	default:
		return 0, ErrInvalidToolsType
	}
}

// DurationOfLastScheduled returns the duration from last time the job successfully scheduled.
func (h *Handler) DurationOfLastScheduled(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		cj, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		if cj.Status.LastScheduleTime == nil {
			return time.Duration(int64(0)), fmt.Errorf("the last time the job was successfully scheduled not found")
		}
		return time.Now().Sub(cj.Status.LastScheduleTime.Time), nil
	case *batchv1.CronJob:
		if val.Status.LastScheduleTime == nil {
			return time.Duration(int64(0)), fmt.Errorf("the last time the job was successfully scheduled not found")
		}
		return time.Now().Sub(val.Status.LastScheduleTime.Time), nil
	case batchv1.CronJob:
		if val.Status.LastScheduleTime == nil {
			return time.Duration(int64(0)), fmt.Errorf("the last time the job was successfully scheduled not found")
		}
		return time.Now().Sub(val.Status.LastScheduleTime.Time), nil
	default:
		return time.Duration(int64(0)), ErrInvalidToolsType
	}
}

// DurationOfCompleted returns the duration from last time the job successful completed.
func (h *Handler) DurationOfCompleted(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		cj, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		if cj.Status.LastSuccessfulTime == nil {
			return time.Duration(int64(0)), fmt.Errorf("the last time the job successfully completed not found")
		}
		return time.Now().Sub(cj.Status.LastSuccessfulTime.Time), nil
	case *batchv1.CronJob:
		if val.Status.LastSuccessfulTime == nil {
			return time.Duration(int64(0)), fmt.Errorf("the last time the job successfully completed not found")
		}
		return time.Now().Sub(val.Status.LastSuccessfulTime.Time), nil
	case batchv1.CronJob:
		if val.Status.LastSuccessfulTime == nil {
			return time.Duration(int64(0)), fmt.Errorf("the last time the job successfully completed not found")
		}
		return time.Now().Sub(val.Status.LastSuccessfulTime.Time), nil
	default:
		return time.Duration(int64(0)), ErrInvalidToolsType
	}
}

// GetSchedule returns the schedule in Cron format.
func (h *Handler) GetSchedule(object interface{}) (string, error) {
	switch val := object.(type) {
	case string:
		cj, err := h.Get(val)
		if err != nil {
			return "", err
		}
		return cj.Spec.Schedule, nil
	case *batchv1.CronJob:
		return val.Spec.Schedule, nil
	case batchv1.CronJob:
		return val.Spec.Schedule, nil
	default:
		return "", ErrInvalidToolsType
	}
}

// IsSuspend check whether the controller will suspend subsequent executions.
func (h *Handler) IsSuspend(object interface{}) (bool, error) {
	// TODO: 当 suspend 没有设置时,是设置成默认值还是返回 "not set" 类型的错误
	switch val := object.(type) {
	case string:
		cj, err := h.Get(val)
		if err != nil {
			return false, err
		}
		// if suspend field not set, set suspend to false.
		if cj.Spec.Suspend == nil {
			return false, nil
		}
		return *cj.Spec.Suspend, nil
	case *batchv1.CronJob:
		if val.Spec.Suspend == nil {
			return false, nil
		}
		return *val.Spec.Suspend, nil
	case batchv1.CronJob:
		if val.Spec.Suspend == nil {
			return false, nil
		}
		return *val.Spec.Suspend, nil
	default:
		return false, ErrInvalidToolsType
	}
}

// GetAge returns cronjob age.
func (h *Handler) GetAge(object interface{}) (time.Duration, error) {
	switch val := object.(type) {
	case string:
		cj, err := h.Get(val)
		if err != nil {
			return time.Duration(int64(0)), err
		}
		return time.Now().Sub(cj.CreationTimestamp.Time), nil
	case *batchv1.CronJob:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	case batchv1.CronJob:
		return time.Now().Sub(val.CreationTimestamp.Time), nil
	default:
		return time.Duration(int64(0)), ErrInvalidToolsType
	}
}

// GetContainers get all container of this cronjob.
func (h *Handler) GetContainers(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		sts, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getContainers(sts), nil
	case *batchv1.CronJob:
		return h.getContainers(val), nil
	case batchv1.CronJob:
		return h.getContainers(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getContainers(sts *batchv1.CronJob) []string {
	var cl []string
	for _, container := range sts.Spec.JobTemplate.Spec.Template.Spec.Containers {
		cl = append(cl, container.Name)
	}
	return cl
}

// GetImages get all container images of this cronjob.
func (h *Handler) GetImages(object interface{}) ([]string, error) {
	switch val := object.(type) {
	case string:
		sts, err := h.Get(val)
		if err != nil {
			return nil, err
		}
		return h.getImages(sts), nil
	case *batchv1.CronJob:
		return h.getImages(val), nil
	case batchv1.CronJob:
		return h.getImages(&val), nil
	default:
		return nil, ErrInvalidToolsType
	}
}
func (h *Handler) getImages(sts *batchv1.CronJob) []string {
	var il []string
	for _, container := range sts.Spec.JobTemplate.Spec.Template.Spec.Containers {
		il = append(il, container.Image)
	}
	return il
}
