package main

import (
	"log"

	"github.com/forbearing/k8s/pod"
)

var (
	namespace = "test"
	filename  = "../../testdata/examples/pod.yaml"
	name      = "mypod"
)

func main() {
	PortForward_Pod()
}

// cleanup will delete or prune created deployments.
func cleanup(handler *pod.Handler) {
	log.Println(handler.Delete(name))
}
