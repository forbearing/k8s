package main

import (
	"log"
	"os"
	"time"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Patch() {
	deployFile := "../../testdata/examples/deployment-patch.yaml"
	patchFile := "../../testdata/examples/deployment-patch-file.yaml"
	deployName := "mydep-patch"
	_, _, _ = deployFile, deployName, patchFile

	handler := deployment.NewOrDie(ctx, "", namespace)
	deploy, err := handler.Apply(deployFile)
	if err != nil {
		log.Fatal(err)
	}
	handler.WaitReady(deployName)
	modifiedDeploy := deploy.DeepCopy()

	// Strategic Merge Patch
	*modifiedDeploy.Spec.Replicas += 1
	deploy2, err := handler.StrategicMergePatch(deploy, modifiedDeploy)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*deploy.Spec.Replicas, *deploy2.Spec.Replicas) //2  3
	handler.WaitReady(unstructName)

	// Json Merge Patch
	*modifiedDeploy.Spec.Replicas += 1
	deploy3, err := handler.MergePatch(deploy, modifiedDeploy)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(*deploy.Spec.Replicas, *deploy3.Spec.Replicas) //2  4
	handler.WaitReady(unstructName)

	// Json Patch
	var patchData []byte
	if patchData, err = os.ReadFile(patchFile); err != nil {
		log.Fatal(err)
	}
	if _, err := handler.JsonPath(deploy, patchData); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)

	// Output:

	//2022/09/05 16:54:11 1 2
	//2022/09/05 16:54:16 1 3
}
