package main

import (
	"log"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/deployment"
)

func main() {
	//Deployment_Create()
	//Deployment_Update()
	Deployment_Update_Status()
	//Deployment_Patch()
	//Deployment_Apply()
	//Deployment_Delete()
	//Deployment_Get()
	//Deployment_List()
	//Deployment_Scale()
	//Deployment_Watch_Single()
	//Deployment_Watch_Label()
	//Deployment_Watch_Field()
	//Deployment_Watch_Namespace()
	//Deployment_Watch_All()
	//Deployment_Informer()
	//Deployment_Tools()
	//Deployment_Others()
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
	k8s.DeleteF(ctx, kubeconfig, filename2, namespace, k8s.IgnoreNotFound)
}
