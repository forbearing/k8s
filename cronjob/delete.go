package cronjob

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Delete deletes cronjob from type string, []byte, *batchv1.CronJob,
// batchv1.CronJob, runtime.Object, *unstructured.Unstructured,
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
	case runtime.Object:
		return h.DeleteFromObject(val)
	default:
		return ErrInvalidDeleteType
	}
}

// DeleteByName deletes cronjob by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.BatchV1().CronJobs(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes cronjob from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes cronjob from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	cmJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	cm := &batchv1.CronJob{}
	if err = json.Unmarshal(cmJson, cm); err != nil {
		return err
	}
	return h.deleteCronjob(cm)
}

// DeleteFromObject deletes cronjob from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	cm, ok := obj.(*batchv1.CronJob)
	if !ok {
		return fmt.Errorf("object type is not *batchv1.CronJob")
	}
	return h.deleteCronjob(cm)
}

// DeleteFromUnstructured deletes cronjob from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	cm := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cm)
	if err != nil {
		return err
	}
	return h.deleteCronjob(cm)
}

// DeleteFromMap deletes cronjob from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	cm := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cm)
	if err != nil {
		return err
	}
	return h.deleteCronjob(cm)
}

// deleteCronjob
func (h *Handler) deleteCronjob(cm *batchv1.CronJob) error {
	var namespace string
	if len(cm.Namespace) != 0 {
		namespace = cm.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.BatchV1().CronJobs(namespace).Delete(h.ctx, cm.Name, h.Options.DeleteOptions)
}
