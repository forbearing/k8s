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

// Get gets job from type string, []byte, *batchv1.Job,
// batchv1.Job, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a job from file path.
func (h *Handler) Get(obj interface{}) (*batchv1.Job, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *batchv1.Job:
		return h.GetFromObject(val)
	case batchv1.Job:
		return h.GetFromObject(&val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	case metav1.Object, runtime.Object:
		return h.GetFromObject(val)
	default:
		return nil, ErrInvalidGetType
	}
}

// GetByName gets job by name.
func (h *Handler) GetByName(name string) (*batchv1.Job, error) {
	return h.clientset.BatchV1().Jobs(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets job from yaml or json file.
func (h *Handler) GetFromFile(filename string) (*batchv1.Job, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets job from bytes data.
func (h *Handler) GetFromBytes(data []byte) (*batchv1.Job, error) {
	jobJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	job := &batchv1.Job{}
	if err = json.Unmarshal(jobJson, job); err != nil {
		return nil, err
	}
	return h.getJob(job)
}

// GetFromObject gets job from metav1.Object or runtime.Object.
func (h *Handler) GetFromObject(obj interface{}) (*batchv1.Job, error) {
	job, ok := obj.(*batchv1.Job)
	if !ok {
		return nil, fmt.Errorf("object type is not *batchv1.Job")
	}
	return h.getJob(job)
}

// GetFromUnstructured gets job from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), job)
	if err != nil {
		return nil, err
	}
	return h.getJob(job)
}

// GetFromMap gets job from unstructured.Unstructured.
func (h *Handler) GetFromMap(u map[string]interface{}) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, job)
	if err != nil {
		return nil, err
	}
	return h.getJob(job)
}

// getJob
// It's necessary to get a new job resource from a old job resource,
// because old job usually don't have job.Status field.
func (h *Handler) getJob(job *batchv1.Job) (*batchv1.Job, error) {
	namespace := job.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.BatchV1().Jobs(namespace).Get(h.ctx, job.Name, h.Options.GetOptions)
}
