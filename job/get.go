package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets job from type string, []byte, *batchv1.Job,
// batchv1.Job, runtime.Object or map[string]interface{}.

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
	case map[string]interface{}:
		return h.GetFromUnstructured(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets job by name.
func (h *Handler) GetByName(name string) (*batchv1.Job, error) {
	return h.clientset.BatchV1().Jobs(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets job from yaml file.
func (h *Handler) GetFromFile(filename string) (*batchv1.Job, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets job from bytes.
func (h *Handler) GetFromBytes(data []byte) (*batchv1.Job, error) {
	jobJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	job := &batchv1.Job{}
	err = json.Unmarshal(jobJson, job)
	if err != nil {
		return nil, err
	}
	return h.getJob(job)
}

// GetFromObject gets job from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*batchv1.Job, error) {
	job, ok := obj.(*batchv1.Job)
	if !ok {
		return nil, fmt.Errorf("object is not *batchv1.Job")
	}
	return h.getJob(job)
}

// GetFromUnstructured gets job from map[string]interface{}.
func (h *Handler) GetFromUnstructured(u map[string]interface{}) (*batchv1.Job, error) {
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
	var namespace string
	if len(job.Namespace) != 0 {
		namespace = job.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.BatchV1().Jobs(namespace).Get(h.ctx, job.Name, h.Options.GetOptions)
}
