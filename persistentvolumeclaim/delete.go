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

// Delete deletes persistentvolumeclaim from type string, []byte, *corev1.PersistentVolumeClaim,
// corev1.PersistentVolumeClaim, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.

// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a persistentvolumeclaim from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *corev1.PersistentVolumeClaim:
		return h.DeleteFromObject(val)
	case corev1.PersistentVolumeClaim:
		return h.DeleteFromObject(&val)
	case runtime.Object:
		return h.DeleteFromObject(val)
	case *unstructured.Unstructured:
		return h.DeleteFromUnstructured(val)
	case unstructured.Unstructured:
		return h.DeleteFromUnstructured(&val)
	case map[string]interface{}:
		return h.DeleteFromMap(val)
	default:
		return ERR_TYPE_DELETE
	}
}

// DeleteByName deletes persistentvolumeclaim by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().PersistentVolumeClaims(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes persistentvolumeclaim from yaml file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes persistentvolumeclaim from bytes.
func (h *Handler) DeleteFromBytes(data []byte) error {
	pvcJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	pvc := &corev1.PersistentVolumeClaim{}
	if err = json.Unmarshal(pvcJson, pvc); err != nil {
		return err
	}
	return h.deletePVC(pvc)
}

// DeleteFromObject deletes persistentvolumeclaim from runtime.Object.
func (h *Handler) DeleteFromObject(obj runtime.Object) error {
	pvc, ok := obj.(*corev1.PersistentVolumeClaim)
	if !ok {
		return fmt.Errorf("object type is not *corev1.PersistentVolumeClaim")
	}
	return h.deletePVC(pvc)
}

// DeleteFromUnstructured deletes persistentvolumeclaim from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	pvc := &corev1.PersistentVolumeClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pvc)
	if err != nil {
		return err
	}
	return h.deletePVC(pvc)
}

// DeleteFromMap deletes persistentvolumeclaim from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	pvc := &corev1.PersistentVolumeClaim{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pvc)
	if err != nil {
		return err
	}
	return h.deletePVC(pvc)
}

// deletePVC
func (h *Handler) deletePVC(pvc *corev1.PersistentVolumeClaim) error {
	var namespace string
	if len(pvc.Namespace) != 0 {
		namespace = pvc.Namespace
	} else {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(h.ctx, pvc.Name, h.Options.DeleteOptions)
}
