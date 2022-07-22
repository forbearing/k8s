package pod

import (
	"io"

	"k8s.io/client-go/tools/remotecommand"
)

type PtyHandler interface {
	io.Reader
	io.Writer
	remotecommand.TerminalSizeQueue
}
