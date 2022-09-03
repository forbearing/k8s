package k8s

import "github.com/forbearing/k8s/dynamic"

var (
	Handler  dynamic.Handler
	New      = dynamic.New
	NewOrDie = dynamic.NewOrDie
)
