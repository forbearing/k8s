package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s/deployment"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testData/examples/deployment.yaml"
	filename2   = "../../testData/examples/deployment-2.yaml"
	update1File = "../../testData/examples/deployment-update1.yaml"
	update2File = "../../testData/examples/deployment-update2.yaml"
	update3File = "../../testData/examples/deployment-update3.yaml"
	nginxFile   = "../../testData/nginx/nginx-deploy.yaml"
	name        = "mydep"
	label       = "type=deployment"
)

var (
	rawName = "mydep-raw"
	rawData = map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name": rawName,
			"labels": map[string]interface{}{
				"app":  rawName,
				"type": "deployment",
			},
		},
		"spec": map[string]interface{}{
			// replicas type is int32, not string.
			"replicas": 1,
			"selector": map[string]interface{}{
				"matchLabels": map[string]interface{}{
					"app":  rawName,
					"type": "deployment",
				},
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app":  rawName,
						"type": "deployment",
					},
				},
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						{
							"name":  "nginx",
							"image": "nginx",
							"resources": map[string]interface{}{
								"limits": map[string]interface{}{
									"cpu": "100m",
								},
							},
						},
					},
				},
			},
		},
	}
)

func main() {
	//Deployment_Create()
	//Deployment_Update()
	//Deployment_Apply()
	//Deployment_Delete()
	//Deployment_Get()
	//Deployment_List()
	//Deployment_Watch()
	Deployment_Tools()
}

func myerr(name string, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success.\n", name)
	}
}

// cleanup will delete or prune created deployments.
func cleanup(handler *deployment.Handler) {
	handler.Delete(name)
	handler.Delete(rawName)
	handler.DeleteFromFile(filename2)
	handler.DeleteFromFile(update1File)
	handler.DeleteFromFile(update2File)
	handler.DeleteFromFile(update3File)
	//k8s.DeleteF(ctx, kubeconfig, nginxFile)
}
