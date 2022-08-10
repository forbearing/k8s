package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../../testdata/examples/deployment.yaml"
	updateFile  = "../../testdata/examples/deployment-update1.yaml"
	filename2   = "../../testdata/nginx/nginx-deploy.yaml"
	name        = "mydep"
	name2       = "nginx-deploy"
	label       = "type=deployment"
)

var (
	unstructName = "mydep-unstruct"
	unstructData = map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name": unstructName,
			"labels": map[string]interface{}{
				"app":  unstructName,
				"type": "deployment",
			},
		},
		"spec": map[string]interface{}{
			// replicas type is int32, not string.
			"replicas": 1,
			"selector": map[string]interface{}{
				"matchLabels": map[string]interface{}{
					"app":  unstructName,
					"type": "deployment",
				},
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"app":  unstructName,
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
	Deployment_Create()
	//Deployment_Update()
	//Deployment_Apply()
	//Deployment_Delete()
	//Deployment_Get()
	//Deployment_List()
	//Deployment_Watch()
	//Deployment_Informer()
	//Deployment_Tools()
	//Deployment_Others()

	// Output:

	//2022/08/03 20:56:05 create deployment from file success: mydep
	//2022/08/03 20:56:05 create deployment from bytes success: mydep
	//2022/08/03 20:56:05 create deployment from *appsv1.Deployment success: mydep
	//2022/08/03 20:56:05 create deployment from appsv1.Deployment success: mydep
	//2022/08/03 20:56:05 create deployment from runtime.Object success: mydep
	//2022/08/03 20:56:05 create deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:56:06 create deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:56:06 create deployment from map[string]interface{} success: mydep-unstruct
	//2022/08/03 20:56:07 update deployment from file success: mydep
	//2022/08/03 20:56:07 update deployment from bytes success: mydep
	//2022/08/03 20:56:07 update deployment from *appsv1.Deployment success: mydep
	//2022/08/03 20:56:07 update deployment from appsv1.Deployment success: mydep
	//2022/08/03 20:56:07 update deployment from runtime.Object success: mydep
	//2022/08/03 20:56:07 update deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:56:07 update deployment from unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:56:07 update deployment from map[string]interface{} success: mydep-unstruct
	//2022/08/03 20:56:08 apply deployment from file (deployment not exists) success: mydep
	//2022/08/03 20:56:08 apply deployment from file (deployment exists) success: mydep
	//2022/08/03 20:56:08 apply deployment from bytes (deployment not exists) success: mydep
	//2022/08/03 20:56:08 apply deployment from bytes (deployment exists) success: mydep
	//2022/08/03 20:56:08 apply deployment from *appsv1.Deployment (deployment not exists) success: mydep
	//2022/08/03 20:56:08 apply deployment from *appsv1.Deployment (deployment exists) success: mydep
	//2022/08/03 20:56:08 apply deployment from appsv1.Deployment (deployment not exists) success: mydep
	//2022/08/03 20:56:09 apply deployment from appsv1.Deployment (deployment exists) success: mydep
	//2022/08/03 20:56:09 apply deployment from runtime.Object (deployment not exists) success: mydep
	//2022/08/03 20:56:10 apply deployment from runtime.Object (deployment exists) success: mydep
	//2022/08/03 20:56:10 apply deployment from *unstructured.Unstructured (deployment not exists) success: mydep-unstruct
	//2022/08/03 20:56:10 apply deployment from *unstructured.Unstructured (deployment exists) success: mydep-unstruct
	//2022/08/03 20:56:11 apply deployment from unstructured.Unstructured (deployment not exists) success: mydep-unstruct
	//2022/08/03 20:56:11 apply deployment from unstructured.Unstructured (deployment exists) success: mydep-unstruct
	//2022/08/03 20:56:12 apply deployment from map[string]interface{} (deployment not exists) success: mydep-unstruct
	//2022/08/03 20:56:12 apply deployment from map[string]interface{} (deployment exists) success: mydep-unstruct
	//2022/08/03 20:56:13 delete deployment by name success:
	//2022/08/03 20:56:13 delete deployment from file success:
	//2022/08/03 20:56:13 delete deployment from bytes success:
	//2022/08/03 20:56:13 delete deployment from *appsv1.Deployment success:
	//2022/08/03 20:56:13 delete deployment from appsv1.Deployment success:
	//2022/08/03 20:56:13 delete deployment from runtime.Object success:
	//2022/08/03 20:56:14 delete deployment from *unstructured.Unstructured success:
	//2022/08/03 20:56:14 delete deployment from unstructured.Unstructured success:
	//2022/08/03 20:56:14 delete deployment from map[string]interface{} success:
	//2022/08/03 20:56:15 get deployment by name success: mydep
	//2022/08/03 20:56:15 get deployment from file success: mydep
	//2022/08/03 20:56:15 get deployment from bytes success: mydep
	//2022/08/03 20:56:15 get deployment from *appsv1.Deployment success: mydep
	//2022/08/03 20:56:15 get deployment from appsv1.Deployment success: mydep
	//2022/08/03 20:56:15 get deployment from runtime.Object success: mydep
	//2022/08/03 20:56:15 get deployment from *unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:56:15 get deployment from unstructured.Unstructured success: mydep-unstruct
	//2022/08/03 20:56:15 get deployment from map[string]interface{} success: mydep-unstruct
	//2022/08/03 20:56:16 ListByLabel success:
	//2022/08/03 20:56:16 [nginx-deploy]
	//2022/08/03 20:56:16 List success:
	//2022/08/03 20:56:16 [nginx-deploy]
	//2022/08/03 20:56:16 ListByNamespace success:
	//2022/08/03 20:56:16 [nginx-deploy]
	//2022/08/03 20:56:16 ListAll success:
	//2022/08/03 20:56:16 [horus-operator ingress-controller calico-kube-controllers coredns metrics-server dashboard-metrics-scraper kubernetes-dashboard local-path-provisioner nfs-provisioner-nfs-subdir-external-provisioner nginx-deploy]
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}

// cleanup will delete or prune created deployments.
func cleanup(handler *deployment.Handler) {
	handler.Delete(name)
	handler.Delete(unstructName)
	handler.DeleteFromFile(updateFile)
	k8s.DeleteF(ctx, kubeconfig, filename2)
}
