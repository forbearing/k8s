package job

import (
	"encoding/json"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create job from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*batchv1.Job, error) {
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

	return h.clientset.BatchV1().Jobs(namespace).Create(h.ctx, job, h.Options.CreateOptions)
}

// CreateFromBytes create job from bytes.
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

	var namespace string
	if len(job.Namespace) != 0 {
		namespace = job.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.BatchV1().Jobs(namespace).Create(h.ctx, job, h.Options.CreateOptions)
}

// CreateFromFile create job from yaml file.
func (h *Handler) CreateFromFile(filename string) (*batchv1.Job, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create job from yaml file, alias to "CreateFromFile".
func (h *Handler) Create(filename string) (*batchv1.Job, error) {
	return h.CreateFromFile(filename)
}
