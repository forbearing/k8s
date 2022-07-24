package job

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Delete deletes job from type string, []byte, *batchv1.Job,
// batchv1.Job, runtime.Object or map[string]interface{}.

// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a job from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *batchv1.Job:
		return h.DeleteFromObject(val)
	case batchv1.Job:
		return h.DeleteFromObject(&val)
	case runtime.Object:
		return h.DeleteFromObject(val)
	case map[string]interface{}:
		return h.DeleteFromUnstructured(val)
	default:
		return ERR_TYPE_DELETE
	}
}

// DeleteByName deletes job by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.BatchV1().Jobs(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes job from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes job from bytes.
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
	return h.deleteJob(job)
}

// DeleteFromObject deletes job from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	job, ok := obj.(*batchv1.Job)
	if !ok {
		return fmt.Errorf("object is not *batchv1.Job")
	}
	return h.deleteJob(job)
}

// DeleteFromUnstructured deletes job from map[string]interface{}.
func (h *Handler) DeleteFromUnstructured(u map[string]interface{}) error {
	job := &batchv1.Job{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, job)
	if err != nil {
		return err
	}
	return h.deleteJob(job)
}

// deleteJob
func (h *Handler) deleteJob(job *batchv1.Job) error {
	var namespace string
	if len(job.Namespace) != 0 {
		namespace = job.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.BatchV1().Jobs(namespace).Delete(h.ctx, job.Name, h.Options.DeleteOptions)
}
