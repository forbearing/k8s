package cronjob

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies cronjob from type string, []byte, *batchv1.CronJob,
// batchv1.CronJob, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Apply(obj interface{}) (*batchv1.CronJob, error) {
	switch val := obj.(type) {
	case string:
		return h.ApplyFromFile(val)
	case []byte:
		return h.ApplyFromBytes(val)
	case *batchv1.CronJob:
		return h.ApplyFromObject(val)
	case batchv1.CronJob:
		return h.ApplyFromObject(&val)
	case runtime.Object:
		return h.ApplyFromObject(val)
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	default:
		return nil, ERR_TYPE_APPLY
	}
}

// ApplyFromFile applies cronjob from yaml file.
func (h *Handler) ApplyFromFile(filename string) (cm *batchv1.CronJob, err error) {
	cm, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if cronjob already exist, update it.
		cm, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply cronjob from bytes.
func (h *Handler) ApplyFromBytes(data []byte) (cm *batchv1.CronJob, err error) {
	cm, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		cm, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies cronjob from runtime.Object.
func (h *Handler) ApplyFromObject(obj runtime.Object) (*batchv1.CronJob, error) {
	cm, ok := obj.(*batchv1.CronJob)
	if !ok {
		return nil, fmt.Errorf("object is not *batchv1.CronJob")
	}
	return h.applyCronjob(cm)
}

// ApplyFromUnstructured applies cronjob from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*batchv1.CronJob, error) {
	cm := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cm)
	if err != nil {
		return nil, err
	}
	return h.applyCronjob(cm)
}

// ApplyFromMap applies cronjob from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*batchv1.CronJob, error) {
	cm := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cm)
	if err != nil {
		return nil, err
	}
	return h.applyCronjob(cm)
}

// applyCronjob
func (h *Handler) applyCronjob(cm *batchv1.CronJob) (*batchv1.CronJob, error) {
	_, err := h.createCronjob(cm)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateCronjob(cm)
	}
	return cm, err
}
