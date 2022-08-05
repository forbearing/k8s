package pod

import (
	"fmt"
	"io"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *corev1.Pod, or corev1.Pod")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *corev1.Pod, corev1.Pod, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ErrInvalidGetType    = ERR_TYPE_CREATE
)

type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}

type LogOptions struct {
	// add '\n' after the string that will be written to "io.Writer"
	NewLine bool

	corev1.PodLogOptions
	io.Writer
}

var DefaultLogOptions = &LogOptions{
	PodLogOptions: corev1.PodLogOptions{},
	Writer:        os.Stdout,
	NewLine:       true,
}
