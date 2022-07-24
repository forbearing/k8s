package job

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies job from type string, []byte, *batchv1.Job,
// batchv1.Job, runtime.Object or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*batchv1.Job, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *batchv1.Job:
		return h.ApplyFromObject(val)
	case batchv1.Job:
		return h.ApplyFromObject(&val)
	case runtime.Object:
		return h.ApplyFromObject(val)
	case map[string]interface{}:
		return h.ApplyFromUnstructured(val)
	default:
		return nil, ERR_TYPE_APPLY
	}
}

// ApplyFromFile applies job from yaml file.
func (h *Handler) ApplyFromFile(filename string) (job *batchv1.Job, err error) {
	job, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if job already exist, update it.
		job, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply job from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (job *batchv1.Job, err error) {
	job, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		job, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies job from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*batchv1.Job, error) {
	job, ok := obj.(*batchv1.Job)
	if !ok {
		return nil, fmt.Errorf("object is not *batchv1.Job")
	}
	return h.applyJob(job)
}

// ApplyFromUnstructured applies job from map[string]interface{}.
func (h *Handler) ApplyFromUnstructured(u map[string]interface{}) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, job)
	if err != nil {
		return nil, err
	}
	return h.applyJob(job)
}

// applyJob
func (h *Handler) applyJob(job *batchv1.Job) (*batchv1.Job, error) {
	_, err := h.createJob(job)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateJob(job)
	}
	return job, err
}
