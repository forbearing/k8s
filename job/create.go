package job

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

// Create creates job from type string, []byte, *batchv1.Job,
// batchv1.Job, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*batchv1.Job, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *batchv1.Job:
		return h.CreateFromObject(val)
	case batchv1.Job:
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
		return nil, ErrInvalidCreateType
	}
}

// CreateFromFile creates job from yaml file.
func (h *Handler) CreateFromFile(filename string) (*batchv1.Job, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates job from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*batchv1.Job, error) {
	jobJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	job := &batchv1.Job{}
	if err = json.Unmarshal(jobJson, job); err != nil {
		return nil, err
	}
	return h.createJob(job)
}

// CreateFromObject creates job from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*batchv1.Job, error) {
	job, ok := obj.(*batchv1.Job)
	if !ok {
		return nil, fmt.Errorf("object type is not *batchv1.Job")
	}
	return h.createJob(job)
}

// CreateFromUnstructured creates job from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), job)
	if err != nil {
		return nil, err
	}
	return h.createJob(job)
}

// CreateFromMap creates job from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, job)
	if err != nil {
		return nil, err
	}
	return h.createJob(job)
}

// createJob
func (h *Handler) createJob(job *batchv1.Job) (*batchv1.Job, error) {
	var namespace string
	if len(job.Namespace) != 0 {
		namespace = job.Namespace
	} else {
		namespace = h.namespace
	}
	job.ResourceVersion = ""
	job.UID = ""
	return h.clientset.BatchV1().Jobs(namespace).Create(h.ctx, job, h.Options.CreateOptions)
}
