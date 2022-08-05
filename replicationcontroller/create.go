package replicationcontroller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates replicationcontroller from type string, []byte,
// *corev1.ReplicationController, corev1.ReplicationController, runtime.Object,
// *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*corev1.ReplicationController, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *corev1.ReplicationController:
		return h.CreateFromObject(val)
	case corev1.ReplicationController:
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

// CreateFromFile creates replicationcontroller from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.ReplicationController, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates replicationcontroller from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.ReplicationController, error) {
	rcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	rc := &corev1.ReplicationController{}
	if err = json.Unmarshal(rcJson, rc); err != nil {
		return nil, err
	}
	return h.createRS(rc)
}

// CreateFromObject creates replicationcontroller from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*corev1.ReplicationController, error) {
	rc, ok := obj.(*corev1.ReplicationController)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.ReplicationController")
	}
	return h.createRS(rc)
}

// CreateFromUnstructured creates replicationcontroller from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*corev1.ReplicationController, error) {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), rc)
	if err != nil {
		return nil, err
	}
	return h.createRS(rc)
}

// CreateFromMap creates replicationcontroller from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*corev1.ReplicationController, error) {
	rc := &corev1.ReplicationController{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, rc)
	if err != nil {
		return nil, err
	}
	return h.createRS(rc)
}

// createRS
func (h *Handler) createRS(rc *corev1.ReplicationController) (*corev1.ReplicationController, error) {
	var namespace string
	if len(rc.Namespace) != 0 {
		namespace = rc.Namespace
	} else {
		namespace = h.namespace
	}
	rc.ResourceVersion = ""
	rc.UID = ""
	return h.clientset.CoreV1().ReplicationControllers(namespace).Create(h.ctx, rc, h.Options.CreateOptions)
}
