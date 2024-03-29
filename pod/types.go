package pod

import (
	"errors"
	"io"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/remotecommand"
)

var (
	ErrInvalidToolsType  = errors.New("type must be string, *corev1.Pod, corev1.Pod, metav1.Object or runtime.Object")
	ErrInvalidCreateType = errors.New("type must be string, []byte, *corev1.Pod, corev1.Pod, metav1.Object, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
	ErrInvalidUpdateType = ErrInvalidCreateType
	ErrInvalidApplyType  = ErrInvalidCreateType
	ErrInvalidDeleteType = ErrInvalidCreateType
	ErrInvalidGetType    = ErrInvalidCreateType
	ErrInvalidLogType    = ErrInvalidCreateType
	ErrInvalidPatchType  = errors.New("patch data type must be string, []byte, *corev1.Pod, corev1.Pod, metav1.Object, runtime.Object, *unstructured.Unstructured, unstructured.Unstructured or map[string]interface{}")
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
