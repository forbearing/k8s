package pod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Delete deletes pod from type string, []byte, *corev1.Pod,
// corev1.Pod, metav1.Object, runtime.Object, *unstructured.Unstructured,
// unstructured.Unstructured or map[string]interface{}.
//
// If passed parameter type is string, it will simply call DeleteByName instead of DeleteFromFile.
// You should always explicitly call DeleteFromFile to delete a pod from file path.
func (h *Handler) Delete(obj interface{}) error {
	switch val := obj.(type) {
	case string:
		return h.DeleteByName(val)
	case []byte:
		return h.DeleteFromBytes(val)
	case *corev1.Pod:
		return h.DeleteFromObject(val)
	case corev1.Pod:
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

// DeleteByName deletes pod by name.
func (h *Handler) DeleteByName(name string) error {
	return h.clientset.CoreV1().Pods(h.namespace).Delete(h.ctx, name, h.Options.DeleteOptions)
}

// DeleteFromFile deletes pod from yaml or json file.
func (h *Handler) DeleteFromFile(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.DeleteFromBytes(data)
}

// DeleteFromBytes deletes pod from bytes data.
func (h *Handler) DeleteFromBytes(data []byte) error {
	podJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	pod := &corev1.Pod{}
	if err = json.Unmarshal(podJson, pod); err != nil {
		return err
	}
	return h.deletePod(pod)
}

// DeleteFromObject deletes pod from metav1.Object or runtime.Object.
func (h *Handler) DeleteFromObject(obj interface{}) error {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("object type is not *corev1.Pod")
	}
	return h.deletePod(pod)
}

// DeleteFromUnstructured deletes pod from *unstructured.Unstructured.
func (h *Handler) DeleteFromUnstructured(u *unstructured.Unstructured) error {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pod)
	if err != nil {
		return err
	}
	return h.deletePod(pod)
}

// DeleteFromMap deletes pod from map[string]interface{}.
func (h *Handler) DeleteFromMap(u map[string]interface{}) error {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pod)
	if err != nil {
		return err
	}
	return h.deletePod(pod)
}

// deletePod
func (h *Handler) deletePod(pod *corev1.Pod) error {
	namespace := pod.GetNamespace()
	if len(namespace) == 0 {
		namespace = h.namespace
	}
	return h.clientset.CoreV1().Pods(namespace).Delete(h.ctx, pod.Name, h.Options.DeleteOptions)
}
