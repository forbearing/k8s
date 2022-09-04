package main

import (
	"encoding/json"
	"io/ioutil"
	"log"

	utilrestmapper "github.com/forbearing/k8s/util/restmapper"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func FindGVK() {
	var (
		err                error
		restMapper         meta.RESTMapper
		gvk                schema.GroupVersionKind
		yamlData, jsonData []byte
		unstructObj        = &unstructured.Unstructured{}
		object             runtime.Object
	)

	if restMapper, err = utilrestmapper.NewRESTMapper(""); err != nil {
		log.Fatal(err)
	}
	if yamlData, err = ioutil.ReadFile(yamlFilename); err != nil {
		log.Fatal(err)
	}
	if jsonData, err = yaml.ToJSON(yamlData); err != nil {
		log.Fatal(err)
	}
	if err = json.Unmarshal(jsonData, unstructObj); err != nil {
		log.Fatal(err)
	}
	object = runtime.Object(unstructObj)

	gvk, err = utilrestmapper.FindGVK(restMapper, yamlFilename)
	checkErr("FindGVK from yaml file", gvk, err)
	gvk, err = utilrestmapper.FindGVK(restMapper, jsonFilename)
	checkErr("FindGVK from json file", gvk, err)

	gvk, err = utilrestmapper.FindGVK(restMapper, yamlData)
	checkErr("FindGVK from yaml data", gvk, err)
	gvk, err = utilrestmapper.FindGVK(restMapper, jsonData)
	checkErr("FindGVK from json data", gvk, err)

	gvk, err = utilrestmapper.FindGVK(restMapper, unstructMap)
	checkErr("FindGVK from map[string]interface{}", gvk, err)
	gvk, err = utilrestmapper.FindGVK(restMapper, unstructObj)
	checkErr("FindGVK from *unstructured.Unstructured", gvk, err)
	gvk, err = utilrestmapper.FindGVK(restMapper, *unstructObj)
	checkErr("FindGVK from unstructured.Unstructured", gvk, err)
	gvk, err = utilrestmapper.FindGVK(restMapper, object)
	checkErr("FindGVK from runtime.Object", gvk, err)
	//gvk, err = utilrestmapper.FindGVK(restMapper, new(int))
	//checkErr("FindGVK from int", gvk, err) // error

	// Output:
	//2022/09/02 14:26:16 FindGVK from yaml file success: apps/v1, Kind=Deployment
	//2022/09/02 14:26:16 FindGVK from json file success: apps/v1, Kind=Deployment
	//2022/09/02 14:26:16 FindGVK from yaml data success: apps/v1, Kind=Deployment
	//2022/09/02 14:26:16 FindGVK from json data success: apps/v1, Kind=Deployment
	//2022/09/02 14:26:16 FindGVK from map[string]interface{} success: apps/v1, Kind=Deployment
	//2022/09/02 14:26:16 FindGVK from *unstructured.Unstructured success: apps/v1, Kind=Deployment
	//2022/09/02 14:26:16 FindGVK from unstructured.Unstructured success: apps/v1, Kind=Deployment
	//2022/09/02 14:26:16 FindGVK from runtime.Object success: apps/v1, Kind=Deployment
	//2022/09/02 14:26:16 FindGVK from int failed: type must be string, []byte, map[string]interface{}, *unstructured.Unstructured, unstructured.Unstructured or runtime.Object
}
