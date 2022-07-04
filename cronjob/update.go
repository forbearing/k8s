package cronjob

import (
	"encoding/json"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// UpdateFromRaw update cronjob from map[string]interface{}.
func (h *Handler) UpdateFromRaw(raw map[string]interface{}) (*batchv1.CronJob, error) {
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

	return h.clientset.BatchV1().CronJobs(namespace).Update(h.ctx, cronjob, h.Options.UpdateOptions)
}

// UpdateFromBytes update cronjob from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*batchv1.CronJob, error) {
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

	return h.clientset.BatchV1().CronJobs(namespace).Update(h.ctx, cronjob, h.Options.UpdateOptions)
}

// UpdateFromFile update cronjob from yaml file.
func (h *Handler) UpdateFromFile(path string) (*batchv1.CronJob, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// Update update cronjob from file, alias to "UpdateFromFile".
func (h *Handler) Update(path string) (*batchv1.CronJob, error) {
	return h.UpdateFromFile(path)
}
