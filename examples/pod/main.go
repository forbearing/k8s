package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/pod"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/pod.yaml"
	filename2   = "../../testdata/nginx/nginx-pod.yaml"
	name        = "mypod"
	name2       = "nginx-pod"
	label       = "type=pod"

	LogPodName = "nginx-logs"
	LogPodData = map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Pod",
		"metadata": map[string]interface{}{
			"name":      LogPodName,
			"namespace": namespace,
			"labels": map[string]interface{}{
				"app":  "LogPodName",
				"type": "pod",
			},
		},
		"spec": map[string]interface{}{
			"containers": []map[string]interface{}{
				{
					"name":  "nginx",
					"image": "nginx",
					"ports": []map[string]interface{}{
						{
							"name":          "http",
							"protocol":      "TCP",
							"containerPort": 80,
						},
						{
							"name":          "https",
							"protocol":      "TCP",
							"containerPort": 443,
						},
					},
				},
			},
		},
	}
)

func main() {
	defer cancel()

	//Pod_Create()
	//Pod_List()
	Pod_Watch()
	//Pod_Tools()
	//Pod_Logs()
	//Pod_Informer()
	//Pod_Others()
}

func cleanup(handler *pod.Handler) {
	handler.Delete(name)
	handler.Delete(LogPodData)
	k8s.DeleteF(ctx, kubeconfig, filename2)
}
func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v.\n", name, val)
	}
}
func wait(handler *pod.Handler, name string) {
	for {
		_, err := handler.Get(name)
		if k8serrors.IsNotFound(err) {
			break
		}
		if err != nil {
			log.Println("handler get pod error: ", err)
			return
		}
		time.Sleep(time.Second)
	}
}
