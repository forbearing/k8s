package job

import (
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply job from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, job)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(job.Namespace) != 0 {
		namespace = job.Namespace
	} else {
		namespace = h.namespace
	}

	_, err = h.clientset.BatchV1().Jobs(namespace).Create(h.ctx, job, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		job, err = h.clientset.BatchV1().Jobs(namespace).Update(h.ctx, job, h.Options.UpdateOptions)
	}
	return job, err
}

// ApplyFromBytes apply job from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (job *batchv1.Job, err error) {
	job, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		job, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply job from yaml file.
func (h *Handler) ApplyFromFile(filename string) (job *batchv1.Job, err error) {
	job, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		job, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply job from yaml file, alias to "ApplyFromFile".
func (h *Handler) Apply(filename string) (*batchv1.Job, error) {
	return h.ApplyFromFile(filename)
}
