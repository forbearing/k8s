package cronjob

import (
	"encoding/json"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// GetFromBytes get cronjob from bytes.
func (h *Handler) GetFromBytes(data []byte) (*batchv1.CronJob, error) {
	cronjobJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cronjob := &batchv1.CronJob{}
	if err = json.Unmarshal(cronjobJson, cronjob); err != nil {
		return nil, err
	}

	var namespace string
	if len(cronjob.Namespace) != 0 {
		namespace = cronjob.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.BatchV1().CronJobs(namespace).Get(h.ctx, cronjob.Name, h.Options.GetOptions)
}

// GetFromFile get cronjob from yaml file.
func (h *Handler) GetFromFile(filename string) (*batchv1.CronJob, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// Get get cronjob by name.
func (h *Handler) Get(name string) (*batchv1.CronJob, error) {
	return h.clientset.BatchV1().CronJobs(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}
