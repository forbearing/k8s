package cronjob

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets cronjob from type string, []byte, *batchv1.CronJob,
// batchv1.CronJob, runtime.Object, *unstructured.Unstructured,
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
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.GetFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.GetFromObject(val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	default:
		return nil, ErrInvalidGetType
	}
}

// GetByName gets cronjob by name.
func (h *Handler) GetByName(name string) (*batchv1.CronJob, error) {
	return h.clientset.BatchV1().CronJobs(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets cronjob from yaml file.
func (h *Handler) GetFromFile(filename string) (*batchv1.CronJob, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets cronjob from bytes.
func (h *Handler) GetFromBytes(data []byte) (*batchv1.CronJob, error) {
	cmJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cm := &batchv1.CronJob{}
	if err = json.Unmarshal(cmJson, cm); err != nil {
		return nil, err
	}
	return h.getCronjob(cm)
}

// GetFromObject gets cronjob from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*batchv1.CronJob, error) {
	cm, ok := obj.(*batchv1.CronJob)
	if !ok {
		return nil, fmt.Errorf("object type is not *batchv1.CronJob")
	}
	return h.getCronjob(cm)
}

// GetFromUnstructured gets cronjob from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*batchv1.CronJob, error) {
	cm := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cm)
	if err != nil {
		return nil, err
	}
	return h.getCronjob(cm)
}

// GetFromMap gets cronjob from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*batchv1.CronJob, error) {
	cm := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cm)
	if err != nil {
		return nil, err
	}
	return h.getCronjob(cm)
}

// getCronjob
// It's necessary to get a new cronjob resource from a old cronjob resource,
// because old cronjob usually don't have cronjob.Status field.
func (h *Handler) getCronjob(cm *batchv1.CronJob) (*batchv1.CronJob, error) {
	var namespace string
	if len(cm.Namespace) != 0 {
		namespace = cm.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.BatchV1().CronJobs(namespace).Get(h.ctx, cm.Name, h.Options.GetOptions)
}
