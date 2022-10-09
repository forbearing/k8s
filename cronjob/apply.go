package cronjob

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// Apply applies cronjob from type string, []byte, *batchv1.CronJob,
// batchv1.CronJob, metav1.Object, runtime.Object, *unstructured.Unstructured,
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
	case *unstructured.Unstructured:
		return h.ApplyFromUnstructured(val)
	case unstructured.Unstructured:
		return h.ApplyFromUnstructured(&val)
	case map[string]interface{}:
		return h.ApplyFromMap(val)
	case metav1.Object, runtime.Object:
		return h.ApplyFromObject(val)
	default:
		return nil, ErrInvalidApplyType
	}
}

// ApplyFromFile applies cronjob from yaml or json file.
func (h *Handler) ApplyFromFile(filename string) (cj *batchv1.CronJob, err error) {
	cj, err = h.CreateFromFile(filename)
	if k8serrors.IsAlreadyExists(err) { // if cronjob already exist, update it.
		cj, err = h.UpdateFromFile(filename)
	}
	return
}

// ApplyFromBytes pply cronjob from bytes data.
func (h *Handler) ApplyFromBytes(data []byte) (cj *batchv1.CronJob, err error) {
	cj, err = h.CreateFromBytes(data)
	if k8serrors.IsAlreadyExists(err) {
		cj, err = h.UpdateFromBytes(data)
	}
	return
}

// ApplyFromObject applies cronjob from metav1.Object or runtime.Object.
func (h *Handler) ApplyFromObject(obj interface{}) (*batchv1.CronJob, error) {
	cj, ok := obj.(*batchv1.CronJob)
	if !ok {
		return nil, fmt.Errorf("object type is not *batchv1.CronJob")
	}
	return h.applyCronjob(cj)
}

// ApplyFromUnstructured applies cronjob from *unstructured.Unstructured.
func (h *Handler) ApplyFromUnstructured(u *unstructured.Unstructured) (*batchv1.CronJob, error) {
	cj := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cj)
	if err != nil {
		return nil, err
	}
	return h.applyCronjob(cj)
}

// ApplyFromMap applies cronjob from map[string]interface{}.
func (h *Handler) ApplyFromMap(u map[string]interface{}) (*batchv1.CronJob, error) {
	cj := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cj)
	if err != nil {
		return nil, err
	}
	return h.applyCronjob(cj)
}

// applyCronjob
func (h *Handler) applyCronjob(cj *batchv1.CronJob) (*batchv1.CronJob, error) {
	_, err := h.createCronjob(cj)
	if k8serrors.IsAlreadyExists(err) {
		return h.updateCronjob(cj)
	}
	return cj, err
}
