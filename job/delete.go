package job

import (
	"encoding/json"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete job from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	jobJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	job := &batchv1.Job{}
	err = json.Unmarshal(jobJson, job)
	if err != nil {
		return err
	}

	var namespace string
	if len(job.Namespace) != 0 {
		namespace = job.Namespace
	} else {
		namespace = h.namespace
	}

	return h.WithNamespace(namespace).DeleteByName(job.Name)
}

// DeleteFromFile delete job from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete job by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.BatchV1().Jobs(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete job by name,alias to "DeleteByName".
func (h *Handler) Delete(name string) error {
	return h.DeleteByName(name)
}
