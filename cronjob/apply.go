package cronjob

import (
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
)

// ApplyFromRaw apply cronjob from map[string]interface{}.
func (h *Handler) ApplyFromRaw(raw map[string]interface{}) (*batchv1.CronJob, error) {
	cronjob := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, cronjob)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(cronjob.Namespace) != 0 {
		namespace = cronjob.Namespace
	} else {
		namespace = h.namespace
	}

	cronjob, err = h.clientset.BatchV1().CronJobs(namespace).Create(h.ctx, cronjob, h.Options.CreateOptions)
	if k8serrors.IsAlreadyExists(err) {
		cronjob, err = h.clientset.BatchV1().CronJobs(namespace).Update(h.ctx, cronjob, h.Options.UpdateOptions)
	}
	return cronjob, err
}

// ApplyFromBytes apply cronjob from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (cronjob *batchv1.CronJob, err error) {
	cronjob, err = h.CreateFromBytes(data)
	if errors.IsAlreadyExists(err) {
		cronjob, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromFile apply cronjob from yaml file
func (h *Handler) ApplyFromFile(filename string) (cronjob *batchv1.CronJob, err error) {
	cronjob, err = h.CreateFromFile(filename)
	if errors.IsAlreadyExists(err) {
		cronjob, err = h.UpdateFromFile(filename)
	}
	return
}

// Apply apply cronjob from file, alias to "ApplyFromFile".
func (h *Handler) Apply(name string) (*batchv1.CronJob, error) {
	return h.ApplyFromFile(name)
}
