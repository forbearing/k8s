package main

import (
	"log"
	"os"

	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/util/signals"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

func Deployment_Patch() {
	var (
		deployFile        = "../../testdata/examples/deployment-patch.yaml"
		strategicYamlFile = "../../testdata/examples/deployment-patch-strategic.yaml"
		strategicJsonFile = "../../testdata/examples/deployment-patch-strategic.json"
		jsontypeYamlFile  = "../../testdata/examples/deployment-patch-jsontype.yaml"
		jsontypeJsonFile  = "../../testdata/examples/deployment-patch-jsontype.json"
		deployName        = "mydep-patch"
	)

	ctx := signals.NewSignalContext()
	handler := deployment.NewOrDie(ctx, "", namespace)
	original, err := handler.Apply(deployFile)
	if err != nil {
		log.Fatal(err)
	}
	modified := original.DeepCopy()
	handler.WaitReady(deployName)

	{
		log.Println("1.1 **JSON Patch** patch data is a filename and the file content is yaml document")
		if _, err := handler.Patch(original, jsontypeYamlFile, types.JSONPatchType); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)
		log.Println("1.2 **JSON Patch** patch data is a filename and the file content is json document")
		if _, err := handler.Patch(original, jsontypeJsonFile, types.JSONPatchType); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)

		log.Println("1.3 **Strategic Merge Patch** patch data is a filename and the file content is yaml document")
		if _, err := handler.Patch(original, strategicYamlFile); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)
		log.Println("1.4 **Strategic Merge Patch** patch data is a filename and the json content is json document")
		if _, err := handler.Patch(original, strategicJsonFile); err != nil {
			log.Fatal(err)
		}

		log.Println("1.5 **JSON Merge Patch** patch data is a filename and the file content is yaml document")
		if _, err := handler.Patch(original, strategicYamlFile, types.MergePatchType); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)
		log.Println("1.6 **JSON Merge Patch** patch data is a filename and the json content is json document")
		if _, err := handler.Patch(original, strategicJsonFile, types.MergePatchType); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)
	}

	{
		log.Println("2.1 **Default to Strategic Merge Patch** patch data is []byte and the content is yaml document")
		var data []byte
		if data, err = os.ReadFile(strategicYamlFile); err != nil {
			log.Fatal(err)
		}
		if _, err := handler.Patch(original, data); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)
		log.Println("2.2 **Default to Strategic Merge Patch** patch data is []byte and the content is json document")
		if data, err = os.ReadFile(strategicJsonFile); err != nil {
			log.Fatal(err)
		}
		if _, err := handler.Patch(original, data); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)
	}

	{

		log.Println("3. **Default to Strategic Merge Patch** patch data is *appsv1.Deployment")
		*modified.Spec.Replicas += 1
		if _, err := handler.Patch(original, modified); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)

		log.Println("4. **Default to Strategic Merge Patch** patch data is appsv1.Deployment")
		*modified.Spec.Replicas += 1
		if _, err := handler.Patch(original, *modified); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)

		log.Println("5. **Default To Strategic Merge Patch** patch data is map[string]interface{}")
		handler.Apply(deployFile)
		unstructMap := make(map[string]interface{})
		if unstructMap, err = runtime.DefaultUnstructuredConverter.ToUnstructured(modified); err != nil {
			log.Fatal(err)
		}
		if _, err := handler.Patch(original, unstructMap); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)

		log.Println("6. **Default to Strategic Merge Patch** patch data is *unstructObj.Unstructured")
		handler.Apply(deployFile)
		unstructObj := &unstructured.Unstructured{Object: unstructMap}
		if _, err := handler.Patch(original, unstructObj); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)

		log.Println("7. **Default to Strategic Merge Patch** patch data is unstructObj.Unstructured")
		handler.Apply(deployFile)
		if _, err := handler.Patch(original, *unstructObj); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)

		log.Println("8.1 **Default to Strategic Merge Patch** patch data is runtime.Object(convert from *appsv1.Deployment)")
		handler.Apply(deployFile)
		object := runtime.Object(modified)
		if _, err := handler.Patch(original, object); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)

		log.Println("8.2 **Default to Strategic Merge Patch** patch data is runtime.Object(convert from *unstructured.Unstructured)")
		handler.Apply(deployFile)
		object = runtime.Object(unstructObj)
		if _, err := handler.Patch(original, object); err != nil {
			log.Fatal(err)
		}
		handler.WaitReady(deployName)
	}

	// Output

	//2022/09/09 14:06:28 1.1 **JSON Patch** patch data is a filename and the file content is yaml document
	//2022/09/09 14:06:31 1.2 **JSON Patch** patch data is a filename and the file content is json document
	//2022/09/09 14:06:34 1.3 **Strategic Merge Patch** patch data is a filename and the file content is yaml document
	//2022/09/09 14:06:40 1.4 **Strategic Merge Patch** patch data is a filename and the json content is json document
	//2022/09/09 14:06:40 1.5 **JSON Merge Patch** patch data is a filename and the file content is yaml document
	//2022/09/09 14:06:48 1.6 **JSON Merge Patch** patch data is a filename and the json content is json document
	//2022/09/09 14:06:48 2.1 **Default to Strategic Merge Patch** patch data is []byte and the content is yaml document
	//2022/09/09 14:06:48 2.2 **Default to Strategic Merge Patch** patch data is []byte and the content is json document
	//2022/09/09 14:06:48 3. **Default to Strategic Merge Patch** patch data is *appsv1.Deployment
	//2022/09/09 14:06:56 4. **Default to Strategic Merge Patch** patch data is appsv1.Deployment
	//2022/09/09 14:07:04 5. **Default To Strategic Merge Patch** patch data is map[string]interface{}
	//2022/09/09 14:07:13 6. **Default to Strategic Merge Patch** patch data is *unstructObj.Unstructured
	//2022/09/09 14:07:13 7. **Default to Strategic Merge Patch** patch data is unstructObj.Unstructured
	//2022/09/09 14:07:19 8.1 **Default to Strategic Merge Patch** patch data is runtime.Object(convert from *appsv1.Deployment)
	//2022/09/09 14:07:19 8.2 **Default to Strategic Merge Patch** patch data is runtime.Object(convert from *unstructured.Unstructured)
}
