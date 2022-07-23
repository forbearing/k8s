package pod

import (
	"fmt"
	"io"

	"k8s.io/client-go/tools/remotecommand"
)

var (
	ERR_TYPE_TOOLS  = fmt.Errorf("type must be string *corev1.Pod, or corev1.Pod")
	ERR_TYPE_CREATE = fmt.Errorf("type must be string, []byte, *corev1.Pod, corev1.Pod, runtime.Object or map[string]interface{}")
	ERR_TYPE_UPDATE = ERR_TYPE_CREATE
	ERR_TYPE_APPLY  = ERR_TYPE_CREATE
	ERR_TYPE_DELETE = ERR_TYPE_CREATE
	ERR_TYPE_GET    = ERR_TYPE_CREATE
)

type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}
