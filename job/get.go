package job

import (
	"encoding/json"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get job from bytes.
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

	var namespace string
	if len(job.Namespace) != 0 {
		namespace = job.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).GetByName(job.Name)
}

// GetFromFile get job from yaml file.
func (h *Handler) GetFromFile(filename string) (*batchv1.Job, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetByName get job by name.
func (h *Handler) GetByName(name string) (*batchv1.Job, error) {
	return h.clientset.BatchV1().Jobs(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// Get get job by name, alias to "GetByName".
func (h *Handler) Get(name string) (*batchv1.Job, error) {
	return h.GetByName(name)
}
