package cronjob

import (
	"encoding/json"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// CreateFromRaw create cronjob from map[string]interface{}.
func (h *Handler) CreateFromRaw(raw map[string]interface{}) (*batchv1.CronJob, error) {
	cronjob := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(raw, cronjob)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(cronjob.Namespace) != 0 {
		namespace = cronjob.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.BatchV1().CronJobs(namespace).Create(h.ctx, cronjob, h.Options.CreateOptions)
}

// CreateFromBytes create cronjob from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*batchv1.CronJob, error) {
	cronjobJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cronjob := &batchv1.CronJob{}
	err = json.Unmarshal(cronjobJson, cronjob)
	if err != nil {
		return nil, err
	}

	var namespace string
	if len(cronjob.Namespace) != 0 {
		namespace = cronjob.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.BatchV1().CronJobs(namespace).Create(h.ctx, cronjob, h.Options.CreateOptions)
}

// CreateFromFile create cronjob from yaml file.
func (h *Handler) CreateFromFile(filename string) (*batchv1.CronJob, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// Create create cronjob from file, alias to "CreateFromFile".
func (h *Handler) Create(path string) (*batchv1.CronJob, error) {
	return h.CreateFromFile(path)
}
