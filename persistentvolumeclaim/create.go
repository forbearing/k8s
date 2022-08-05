package persistentvolumeclaim

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

// Create creates persistentvolumeclaim from type string, []byte, *corev1.PersistentVolumeClaim,
// corev1.PersistentVolumeClaim, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*corev1.PersistentVolumeClaim, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *corev1.PersistentVolumeClaim:
		return h.CreateFromObject(val)
	case corev1.PersistentVolumeClaim:
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

// CreateFromFile creates persistentvolumeclaim from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.PersistentVolumeClaim, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates persistentvolumeclaim from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.PersistentVolumeClaim, error) {
	pvcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pvc := &corev1.PersistentVolumeClaim{}
	if err = json.Unmarshal(pvcJson, pvc); err != nil {
		return nil, err
	}
	return h.createPVC(pvc)
}

// CreateFromObject creates persistentvolumeclaim from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*corev1.PersistentVolumeClaim, error) {
	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.PersistentVolumeClaim")
	}
	return h.createPVC(pvc)
}

// CreateFromUnstructured creates persistentvolumeclaim from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pvc)
	if err != nil {
		return nil, err
	}
	return h.createPVC(pvc)
}

// CreateFromMap creates persistentvolumeclaim from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pvc)
	if err != nil {
		return nil, err
	}
	return h.createPVC(pvc)
}

// createPVC
func (h *Handler) createPVC(pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	var namespace string
	if len(pvc.Namespace) != 0 {
		namespace = pvc.Namespace
	} else {
		namespace = h.namespace
	}
	pvc.ResourceVersion = ""
	pvc.UID = ""
	return h.clientset.CoreV1().PersistentVolumeClaims(namespace).Create(h.ctx, pvc, h.Options.CreateOptions)
}
