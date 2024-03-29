package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Update updates job from type string, []byte, *batchv1.Job,
// batchv1.Job, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*batchv1.Job, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *batchv1.Job:
		return h.UpdateFromObject(val)
	case batchv1.Job:
		return h.UpdateFromObject(&val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	case metav1.Object, runtime.Object:
		return h.UpdateFromObject(val)
	default:
		return nil, ErrInvalidUpdateType
	}
}

// UpdateFromFile updates job from yaml or json file.
func (h *Handler) UpdateFromFile(filename string) (*batchv1.Job, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates job from bytes data.
func (h *Handler) UpdateFromBytes(data []byte) (*batchv1.Job, error) {
	jobJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	job := &batchv1.Job{}
	if err = json.Unmarshal(jobJson, job); err != nil {
		return nil, err
	}
	return h.updateJob(job)
}

// UpdateFromObject updates job from metav1.Object or runtime.Object.
func (h *Handler) UpdateFromObject(obj interface{}) (*batchv1.Job, error) {
	job, ok := obj.(*batchv1.Job)
	if !ok {
		return nil, fmt.Errorf("object type is not *batchv1.Job")
	}
	return h.updateJob(job)
}

// UpdateFromUnstructured updates job from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), job)
	if err != nil {
		return nil, err
	}
	return h.updateJob(job)
}

// UpdateFromMap updates job from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, job)
	if err != nil {
		return nil, err
	}
	return h.updateJob(job)
}

// updateJob
func (h *Handler) updateJob(job *batchv1.Job) (*batchv1.Job, error) {
	namespace := job.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	//// resourceVersion cann't be set, the resourceVersion field is empty.
	job.ResourceVersion = ""
	job.UID = ""
	return h.clientset.BatchV1().Jobs(namespace).Update(h.ctx, job, h.Options.UpdateOptions)
}
