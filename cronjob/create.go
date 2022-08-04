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

// Create creates cronjob from type string, []byte, *batchv1.CronJob,
// batchv1.CronJob, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*batchv1.CronJob, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *batchv1.CronJob:
		return h.CreateFromObject(val)
	case batchv1.CronJob:
		return h.CreateFromObject(&val)
	case runtime.Object:
		if reflect.TypeOf(val).String() == "*unstructured.Unstructured" {
			return h.CreateFromUnstructured(val.(*unstructured.Unstructured))
		}
		return h.CreateFromObject(val)
	case *unstructured.Unstructured:
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates cronjob from yaml file.
func (h *Handler) CreateFromFile(filename string) (*batchv1.CronJob, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates cronjob from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*batchv1.CronJob, error) {
	cjJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	cj := &batchv1.CronJob{}
	if err = json.Unmarshal(cjJson, cj); err != nil {
		return nil, err
	}
	return h.createCronjob(cj)
}

// CreateFromObject creates cronjob from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*batchv1.CronJob, error) {
	cj, ok := obj.(*batchv1.CronJob)
	if !ok {
		return nil, fmt.Errorf("object type is not *batchv1.CronJob")
	}
	return h.createCronjob(cj)
}

// CreateFromUnstructured creates cronjob from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*batchv1.CronJob, error) {
	cj := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), cj)
	if err != nil {
		return nil, err
	}
	return h.createCronjob(cj)
}

// CreateFromMap creates cronjob from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*batchv1.CronJob, error) {
	cj := &batchv1.CronJob{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, cj)
	if err != nil {
		return nil, err
	}
	return h.createCronjob(cj)
}

// createCronjob
func (h *Handler) createCronjob(cj *batchv1.CronJob) (*batchv1.CronJob, error) {
	var namespace string
	if len(cj.Namespace) != 0 {
		namespace = cj.Namespace
	} else {
		namespace = h.namespace
	}
	cj.ResourceVersion = ""
	cj.UID = ""
	return h.clientset.BatchV1().CronJobs(namespace).Create(h.ctx, cj, h.Options.CreateOptions)
}
