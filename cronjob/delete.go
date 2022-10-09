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

// Delete deletes cronjob from type string, []byte, *batchv1.CronJob,
// batchv1.CronJob, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a cronjob from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *batchv1.CronJob:
		return h.DeleteFromObject(val)
	case batchv1.CronJob:
		return h.DeleteFromObject(&val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	case metav1.Object, runtime.Object:
		return h.DeleteFromObject(val)
	default:
		return ErrInvalidDeleteType
	}
}

// DeleteByName deletes cronjob by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.BatchV1().CronJobs(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes cronjob from yaml or json file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes cronjob from bytes data.
func (h *Handler) DeleteFromBytes(data []byte) error {
	cmJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	cj := &batchv1.CronJob{}
	if err = json.Unmarshal(cmJson, cj); err != nil {
		return err
	}
	return h.deleteCronjob(cj)
}

// DeleteFromObject deletes cronjob from metav1.Object or runtime.Object.
func (h *Handler) DeleteFromObject(obj interface{}) error {
	cj, ok := obj.(*batchv1.CronJob)
	if !ok {
		return fmt.Errorf("object type is not *batchv1.CronJob")
	}
	return h.deleteCronjob(cj)
}

// DeleteFromUnstructured deletes cronjob from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	cj := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cj)
	if err != nil {
		return err
	}
	return h.deleteCronjob(cj)
}

// DeleteFromMap deletes cronjob from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	cj := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cj)
	if err != nil {
		return err
	}
	return h.deleteCronjob(cj)
}

// deleteCronjob
func (h *Handler) deleteCronjob(cj *batchv1.CronJob) error {
	namespace := cj.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.BatchV1().CronJobs(namespace).Delete(h.ctx, cj.Name, h.Options.DeleteOptions)
}
