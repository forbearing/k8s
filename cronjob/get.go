package cronjob

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets cronjob from type string, []byte, *batchv1.CronJob,
// batchv1.CronJob, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a cronjob from file path.
func (h *Handler) Get(obj interface{}) (*batchv1.CronJob, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *batchv1.CronJob:
		return h.GetFromObject(val)
	case batchv1.CronJob:
		return h.GetFromObject(&val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	case metav1.Object, runtime.Object:
		return h.GetFromObject(val)
	default:
		return nil, ErrInvalidGetType
	}
}

// GetByName gets cronjob by name.
func (h *Handler) GetByName(name string) (*batchv1.CronJob, error) {
	return h.clientset.BatchV1().CronJobs(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets cronjob from yaml or json file.
func (h *Handler) GetFromFile(filename string) (*batchv1.CronJob, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets cronjob from bytes data.
func (h *Handler) GetFromBytes(data []byte) (*batchv1.CronJob, error) {
	cmJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cj := &batchv1.CronJob{}
	if err = json.Unmarshal(cmJson, cj); err != nil {
		return nil, err
	}
	return h.getCronjob(cj)
}

// GetFromObject gets cronjob from metav1.Object or runtime.Object.
func (h *Handler) GetFromObject(obj interface{}) (*batchv1.CronJob, error) {
	cj, ok := obj.(*batchv1.CronJob)
	if !ok {
		return nil, fmt.Errorf("object type is not *batchv1.CronJob")
	}
	return h.getCronjob(cj)
}

// GetFromUnstructured gets cronjob from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*batchv1.CronJob, error) {
	cj := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cj)
	if err != nil {
		return nil, err
	}
	return h.getCronjob(cj)
}

// GetFromMap gets cronjob from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*batchv1.CronJob, error) {
	cj := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cj)
	if err != nil {
		return nil, err
	}
	return h.getCronjob(cj)
}

// getCronjob
// It's necessary to get a new cronjob resource from a old cronjob resource,
// because old cronjob usually don't have cronjob.Status field.
func (h *Handler) getCronjob(cj *batchv1.CronJob) (*batchv1.CronJob, error) {
	namespace := cj.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.BatchV1().CronJobs(namespace).Get(h.ctx, cj.Name, h.Options.GetOptions)
}
