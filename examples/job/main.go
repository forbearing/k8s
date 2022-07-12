package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s/job"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/job.yaml"
	filename2   = "../../testdata/examples/job-failed.yaml"
	name        = "myjob"
	name2       = "myjob-failed"
	label       = "type=job"
)

func main() {
	Job_Tools()
}

func checkErr(name string, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success.\n", name)
	}
}
func cleanup(handler *job.Handler) {
	handler.Delete(name)
	handler.Delete(name2)
}
