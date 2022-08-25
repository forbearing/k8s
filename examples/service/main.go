package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s/service"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/service.yaml"
	filenameNP  = "../../testdata/examples/service-nodeport.yaml"
	filenameLB  = "../../testdata/examples/service-loadbalancer.yaml"
	filenameEN  = "../../testdata/examples/service-externalname.yaml"
	filenameEI  = "../../testdata/examples/service-externalips.yaml"
	name        = "mysvc"
	nameNP      = "mysvc-nodeport"
	nameLB      = "mysvc-loadbalancer"
	nameEN      = "mysvc-externalname"
	nameEI      = "mysvc-externalip"
	label       = "type=service"
)

func main() {
	//Service_Tools()
	Service_Informer()
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}
func cleanup(handler *service.Handler) {
	handler.Delete(name)
	handler.Delete(nameNP)
	handler.Delete(nameLB)
	handler.Delete(nameEN)
	handler.Delete(nameEI)
}
