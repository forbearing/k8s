package job

import (
	"encoding/json"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update job from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*batchv1.Job, error) {
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

	return h.clientset.BatchV1().Jobs(namespace).Update(h.ctx, job, h.Options.UpdateOptions)
}

// UpdateFromBytes update job from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*batchv1.Job, error) {
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

	return h.clientset.BatchV1().Jobs(namespace).Update(h.ctx, job, h.Options.UpdateOptions)
}

// UpdateFromFile update job from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*batchv1.Job, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update job from yaml file, alias to "UpdateFromFile".
func (h *Handler) Update(filename string) (*batchv1.Job, error) {
	return h.UpdateFromFile(filename)
}
