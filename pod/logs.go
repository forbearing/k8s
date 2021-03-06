package pod

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Log get pod logs from type string, []byte, *corev1.pod, corev1.pod,
// runtime.Object, *unstructured.Unstructured, unstructured.Unstructured
// or map[string]interface{}.

// If passed parameter type is string, it will simply call LogByName instead of LogFromFile.
// You should always explicitly call LogFromFile to get pod logs from file path.
func (h *Handler) Log(obj interface{}, logOptions *LogOptions) error {
	switch val := obj.(type) {
	case string:
		return h.LogByName(val, logOptions)
	case []byte:
		return h.LogFromBytes(val, logOptions)
	case *corev1.Pod:
		return h.LogFromObject(val, logOptions)
	case corev1.Pod:
		return h.LogFromObject(&val, logOptions)
	case runtime.Object:
		return h.LogFromObject(val, logOptions)
	case *unstructured.Unstructured:
		return h.LogFromUnstructured(val, logOptions)
	case unstructured.Unstructured:
		return h.LogFromUnstructured(&val, logOptions)
	case map[string]interface{}:
		return h.LogFromMap(val, logOptions)
	default:
		return ERR_TYPE_CREATE
	}
}

// LogByName gets pod by name.
func (h *Handler) LogByName(name string, logOption *LogOptions) error {
	return h.getLog(h.namespace, name, logOption)
}

// LogFromFile get pod logs from yaml file.
func (h *Handler) LogFromFile(filename string, logOptions *LogOptions) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return h.LogFromBytes(data, logOptions)
}

// LogFromBytes get pod logs from bytes.
func (h *Handler) LogFromBytes(data []byte, logOptions *LogOptions) error {
	podJson, err := yaml.ToJSON(data)
	if err != nil {
		return err
	}

	pod := &corev1.Pod{}
	if err = json.Unmarshal(podJson, pod); err != nil {
		return err
	}
	return h.logPod(pod, logOptions)
}

// LogFromObject get logs from runtime.Object.
func (h *Handler) LogFromObject(obj runtime.Object, logOptions *LogOptions) error {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return fmt.Errorf("object type is not *corev1.Pod")
	}
	return h.logPod(pod, logOptions)
}

// LogFromUnstructured get logs from *unstructured.Unstructured.
func (h *Handler) LogFromUnstructured(u *unstructured.Unstructured, logOptions *LogOptions) error {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u.UnstructuredContent(), pod)
	if err != nil {
		return err
	}
	return h.logPod(pod, logOptions)
}

// LogFromMap get logs from map[string]interface{}.
func (h *Handler) LogFromMap(u map[string]interface{}, logOptions *LogOptions) error {
	pod := &corev1.Pod{}
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(u, pod)
	if err != nil {
		return err
	}
	return h.logPod(pod, logOptions)
}

// logPod
func (h *Handler) logPod(pod *corev1.Pod, logOptions *LogOptions) error {
	if !h.IsReady(pod.Name) {
		return fmt.Errorf("pod/%s is not ready", pod.Name)
	}

	var namespace string
	if len(pod.Namespace) != 0 {
		namespace = pod.Namespace
	} else {
		namespace = h.namespace
	}
	return h.getLog(namespace, pod.Name, logOptions)
}

func (h *Handler) getLog(namespace, name string, logOptions *LogOptions) error {
	req := h.clientset.CoreV1().Pods(namespace).GetLogs(name, &logOptions.PodLogOptions)
	readCloser, err := req.Stream(h.ctx)
	if err != nil {
		return err
	}
	defer readCloser.Close()

	scanner := bufio.NewScanner(readCloser)
	scanner.Split(bufio.ScanLines)
	if !logOptions.NewLine {
		for scanner.Scan() {
			fmt.Fprintf(logOptions.Writer, "%s", scanner.Text())
		}
		return scanner.Err()
	}
	for scanner.Scan() {
		fmt.Fprintf(logOptions.Writer, "%s\n", scanner.Text())
	}
	return scanner.Err()
}
