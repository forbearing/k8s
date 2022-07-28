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

// Update updates persistentvolumeclaim from type string, []byte, *corev1.PersistentVolumeClaim,
// corev1.PersistentVolumeClaim, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
func (h *Handler) Update(obj interface{}) (*corev1.PersistentVolumeClaim, error) {
	switch val := obj.(type) {
	case string:
		return h.UpdateFromFile(val)
	case []byte:
		return h.UpdateFromBytes(val)
	case *corev1.PersistentVolumeClaim:
		return h.UpdateFromObject(val)
	case corev1.PersistentVolumeClaim:
		return h.UpdateFromObject(&val)
	case runtime.Object:
		return h.UpdateFromObject(val)
	case *unstructured.Unstructured:
		return h.UpdateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.UpdateFromUnstructured(&val)
	case map[string]interface{}:
		return h.UpdateFromMap(val)
	default:
		return nil, ERR_TYPE_UPDATE
	}
}

// UpdateFromFile updates persistentvolumeclaim from yaml file.
func (h *Handler) UpdateFromFile(filename string) (*corev1.PersistentVolumeClaim, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.UpdateFromBytes(data)
}

// UpdateFromBytes updates persistentvolumeclaim from bytes.
func (h *Handler) UpdateFromBytes(data []byte) (*corev1.PersistentVolumeClaim, error) {
	pvcJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pvc := &corev1.PersistentVolumeClaim{}
	if err = json.Unmarshal(pvcJson, pvc); err != nil {
		return nil, err
	}
	return h.updatePVC(pvc)
}

// UpdateFromObject updates persistentvolumeclaim from runtime.Object.
func (h *Handler) UpdateFromObject(obj runtime.Object) (*corev1.PersistentVolumeClaim, error) {
	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.PersistentVolumeClaim")
	}
	return h.updatePVC(pvc)
}

// UpdateFromUnstructured updates persistentvolumeclaim from *unstructured.Unstructured.
func (h *Handler) UpdateFromUnstructured(u *unstructured.Unstructured) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pvc)
	if err != nil {
		return nil, err
	}
	return h.updatePVC(pvc)
}

// UpdateFromMap updates persistentvolumeclaim from map[string]interface{}.
func (h *Handler) UpdateFromMap(u map[string]interface{}) (*corev1.PersistentVolumeClaim, error) {
	pvc := &corev1.PersistentVolumeClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pvc)
	if err != nil {
		return nil, err
	}
	return h.updatePVC(pvc)
}

// updatePVC
func (h *Handler) updatePVC(pvc *corev1.PersistentVolumeClaim) (*corev1.PersistentVolumeClaim, error) {
	var namespace string
	if len(pvc.Namespace) != 0 {
		namespace = pvc.Namespace
	} else {
		namespace = h.namespace
	}
	pvc.ResourceVersion = ""
	pvc.UID = ""
	return h.clientset.CoreV1().PersistentVolumeClaims(namespace).Update(h.ctx, pvc, h.Options.UpdateOptions)
}
