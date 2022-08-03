package pod

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Create creates pod from type string, []byte, *corev1.pod, corev1.pod,
// runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.
func (h *Handler) Create(obj interface{}) (*corev1.Pod, error) {
	switch val := obj.(type) {
	case string:
		return h.CreateFromFile(val)
	case []byte:
		return h.CreateFromBytes(val)
	case *corev1.Pod:
		return h.CreateFromObject(val)
	case corev1.Pod:
		return h.CreateFromObject(&val)
	case runtime.Object:
		return h.CreateFromObject(val)
	case *unstructured.Unstructured:
		// - 如果传入的类型是 *unstructured.Unstructured 做类型断言时,它会自动转换成
		//   runtime.Object 类型, 而不是 *unstructured.Unstructured
		// - 所以不支持从 *unstructured.Unstructured 来创建 pod
		//   只支持从 unstructured.Unstructured 来创建 pod
		return h.CreateFromUnstructured(val)
	case unstructured.Unstructured:
		return h.CreateFromUnstructured(&val)
	case map[string]interface{}:
		return h.CreateFromMap(val)
	default:
		return nil, ERR_TYPE_CREATE
	}
}

// CreateFromFile creates pod from yaml file.
func (h *Handler) CreateFromFile(filename string) (*corev1.Pod, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return h.CreateFromBytes(data)
}

// CreateFromBytes creates pod from bytes.
func (h *Handler) CreateFromBytes(data []byte) (*corev1.Pod, error) {
	podJson, err := yaml.ToJSON(data)
	if err != nil {
		return nil, err
	}

	pod := &corev1.Pod{}
	if err = json.Unmarshal(podJson, pod); err != nil {
		return nil, err
	}
	return h.createPod(pod)
}

// CreateFromObject creates pod from runtime.Object.
func (h *Handler) CreateFromObject(obj runtime.Object) (*corev1.Pod, error) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return nil, fmt.Errorf("object type is not *corev1.Pod")
	}
	return h.createPod(pod)
}

// CreateFromUnstructured creates pod from *unstructured.Unstructured.
func (h *Handler) CreateFromUnstructured(u *unstructured.Unstructured) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pod)
	if err != nil {
		return nil, err
	}
	return h.createPod(pod)
}

// CreateFromMap creates pod from map[string]interface{}.
func (h *Handler) CreateFromMap(u map[string]interface{}) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pod)
	if err != nil {
		return nil, err
	}
	return h.createPod(pod)
}

// createPod
func (h *Handler) createPod(pod *corev1.Pod) (*corev1.Pod, error) {
	var namespace string
	if len(pod.Namespace) != 0 {
		namespace = pod.Namespace
	} else {
		namespace = h.namespace
	}
	pod.UID = ""
	pod.ResourceVersion = ""
	return h.clientset.CoreV1().Pods(namespace).Create(h.ctx, pod, h.Options.CreateOptions)
}
