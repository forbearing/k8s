package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates job from type string, []byte, *batchv1.Job,
// batchv1.Job, runtime.Object or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*batchv1.Job, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *batchv1.Job:
		return h.CreateFromObject(val)
	case batchv1.Job:
		return h.CreateFromObject(&val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case map[string]interface{}:
		return h.CreateFromUnstructured(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates job from yaml file.
func (h *Handler) CreateFromFile(filename string) (*batchv1.Job, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates job from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*batchv1.Job, error) {
	jobJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	job := &batchv1.Job{}
	err = json.Unmarshal(jobJson, job)
	if err != nil {
		return nil, err
	}
	return h.createJob(job)
}

// CreateFromObject creates job from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*batchv1.Job, error) {
	job, ok := obj.(*batchv1.Job)
	if !ok {
		return nil, fmt.Errorf("object is not *batchv1.Job")
	}
	return h.createJob(job)
}

// CreateFromUnstructured creates job from map[string]interface{}.
func (h *Handler) CreateFromUnstructured(u map[string]interface{}) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, job)
	if err != nil {
		return nil, err
	}
	return h.createJob(job)
}

// createJob
func (h *Handler) createJob(job *batchv1.Job) (*batchv1.Job, error) {
	var namespace string
	if len(job.Namespace) != 0 {
		namespace = job.Namespace
	} else {
		namespace = h.namespace
	}
	job.ResourceVersion = ""
	job.UID = ""
	return h.clientset.BatchV1().Jobs(namespace).Create(h.ctx, job, h.Options.CreateOptions)
}
