package cronjob

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates cronjob from type string, []byte, *batchv1.CronJob,
// batchv1.CronJob, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*batchv1.CronJob, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *batchv1.CronJob:
		return h.UpdateFromObject(val)
	case batchv1.CronJob:
		return h.UpdateFromObject(&val)
	case runtime.Object:
		return h.UpdateFromObject(val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	default:
		return nil, ERR_TYPE_UPDATE
	}
}

// UpdateFromFile updates cronjob from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*batchv1.CronJob, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates cronjob from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*batchv1.CronJob, error) {
	cjJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cj := &batchv1.CronJob{}
	err = json.Unmarshal(cjJson, cj)
	if err != nil {
		return nil, err
	}
	return h.updateCronjob(cj)
}

// UpdateFromObject updates cronjob from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*batchv1.CronJob, error) {
	cj, ok := obj.(*batchv1.CronJob)
	if !ok {
		return nil, fmt.Errorf("object is not *batchv1.CronJob")
	}
	return h.updateCronjob(cj)
}

// UpdateFromUnstructured updates cronjob from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*batchv1.CronJob, error) {
	cj := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cj)
	if err != nil {
		return nil, err
	}
	return h.updateCronjob(cj)
}

// UpdateFromMap updates cronjob from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*batchv1.CronJob, error) {
	cj := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cj)
	if err != nil {
		return nil, err
	}
	return h.updateCronjob(cj)
}

// updateCronjob
func (h *Handler) updateCronjob(cj *batchv1.CronJob) (*batchv1.CronJob, error) {
	var namespace string
	if len(cj.Namespace) != 0 {
		namespace = cj.Namespace
	} else {
		namespace = h.namespace
	}
	//// resourceVersion cann't be set, the resourceVersion field is empty.
	cj.ResourceVersion = ""
	cj.UID = ""
	return h.clientset.BatchV1().CronJobs(namespace).Update(h.ctx, cj, h.Options.UpdateOptions)
}
