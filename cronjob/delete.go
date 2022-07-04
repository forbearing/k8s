package cronjob

import (
	"encoding/json"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// DeleteFromBytes delete cronjob from bytes.
func (h *Handler) DeleteFromBytes(data []byte) (err error) {
	cronjobJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	cronjob := &batchv1.CronJob{}
	err = json.Unmarshal(cronjobJson, cronjob)
	if err != nil {
		return err
	}

	var namespace string
	if len(cronjob.Namespace) != 0 {
		namespace = cronjob.Namespace
	} else {
		namespace = h.namespace
	}

	return h.clientset.BatchV1().CronJobs(namespace).Delete(h.ctx, cronjob.Name, h.Options.DeleteOptions)
}

// DeleteFromFile delete cronjob from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteByName delete cronjob by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.BatchV1().CronJobs(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// Delete delete cronjob by name, alias to "DeleteByName".
func (h *Handler) Delete(name string) (err error) {
	return h.DeleteByName(name)
}
