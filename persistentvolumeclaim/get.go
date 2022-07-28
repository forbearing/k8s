package persistentvolumeclaim

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Get gets persistentvolumeclaim from type string, []byte, *corev1.PersistentVolumeClaim,
// corev1.PersistentVolumeClaim, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call GetByName instead of GetFromFile.
// You should always explicitly call GetFromFile to get a persistentvolumeclaim from file path.
func (h *Handler) Get(obj interface{}) (*corev1.PersistentVolumeClaim, error) {
	switch val := obj.(type) {
	case string:
		return h.GetByName(val)
	case []byte:
		return h.GetFromBytes(val)
	case *corev1.PersistentVolumeClaim:
		return h.GetFromObject(val)
	case corev1.PersistentVolumeClaim:
		return h.GetFromObject(&val)
	case runtime.Object:
		return h.GetFromObject(val)
	case *unstructured.Unstructured:
		return h.GetFromUnstructured(val)
	case unstructured.Unstructured:
		return h.GetFromUnstructured(&val)
	case map[string]interface{}:
		return h.GetFromMap(val)
	default:
		return nil, ERR_TYPE_GET
	}
}

// GetByName gets persistentvolumeclaim by name.
func (h *Handler) GetByName(name string) (*corev1.PersistentVolumeClaim, error) {
	return h.clientset.CoreV1().PersistentVolumeClaims(h.namespace).Get(h.ctx, name, h.Options.GetOptions)
}

// GetFromFile gets persistentvolumeclaim from yaml file.
func (h *Handler) GetFromFile(filename string) (*corev1.PersistentVolumeClaim, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.GetFromBytes(data)
}

// GetFromBytes gets persistentvolumeclaim from bytes.
func (h *Handler) GetFromBytes(data []byte) (*corev1.PersistentVolumeClaim, error) {
	pvcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pvc := &corev1.PersistentVolumeClaim{}
	err = json.Unmarshal(pvcJson, pvc)
	if err != nil {
		return nil, err
	}
	return h.getPVC(pvc)
}

// GetFromObject gets persistentvolumeclaim from runtime.Object.
func (h *Handler) GetFromObject(obj runtime.Object) (*corev1.PersistentVolumeClaim, error) {
	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		return nil, fmt.Errorf("object is not *corev1.PersistentVolumeClaim")
	}
	return h.getPVC(pvc)
}

// GetFromUnstructured gets persistentvolumeclaim from *unstructured.Unstructured.
func (h *Handler) GetFromUnstructured(u *unstructured.Unstructured) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pvc)
	if err != nil {
		return nil, err
	}
	return h.getPVC(pvc)
}

// GetFromMap gets persistentvolumeclaim from map[string]interface{}.
func (h *Handler) GetFromMap(u map[string]interface{}) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pvc)
	if err != nil {
		return nil, err
	}
	return h.getPVC(pvc)
}

// getPVC
// It's necessary to get a new persistentvolumeclaim resource from a old persistentvolumeclaim resource,
// because old persistentvolumeclaim usually don't have persistentvolumeclaim.Status field.
func (h *Handler) getPVC(pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	var namespace string
	if len(pvc.Namespace) != 0 {
		namespace = pvc.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().PersistentVolumeClaims(namespace).Get(h.ctx, pvc.Name, h.Options.GetOptions)
}
