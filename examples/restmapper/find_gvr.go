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

func FindGVR() {
	var (
		err                error
		restMapper         meta.RESTMapper
		gvr                schema.GroupVersionResource
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

	gvr, err = utilrestmapper.FindGVR(restMapper, yamlFilename)
	checkErr("FindGVR from yaml file", gvr, err)
	gvr, err = utilrestmapper.FindGVR(restMapper, jsonFilename)
	checkErr("FindGVR from json file", gvr, err)

	gvr, err = utilrestmapper.FindGVR(restMapper, yamlData)
	checkErr("FindGVR from yaml data", gvr, err)
	gvr, err = utilrestmapper.FindGVR(restMapper, jsonData)
	checkErr("FindGVR from json data", gvr, err)

	gvr, err = utilrestmapper.FindGVR(restMapper, unstructMap)
	checkErr("FindGVR from map[string]interface{}", gvr, err)
	gvr, err = utilrestmapper.FindGVR(restMapper, unstructObj)
	checkErr("FindGVR from *unstructured.Unstructured", gvr, err)
	gvr, err = utilrestmapper.FindGVR(restMapper, *unstructObj)
	checkErr("FindGVR from unstructured.Unstructured", gvr, err)
	gvr, err = utilrestmapper.FindGVR(restMapper, object)
	checkErr("FindGVR from runtime.Object", gvr, err)
	// error
	gvr, err = utilrestmapper.FindGVR(restMapper, new(int))
	checkErr("FindGVR from int", gvr, err)

	// Output:
	//2022/09/02 14:30:13 FindGVR from yaml file success: apps/v1, Resource=deployments
	//2022/09/02 14:30:13 FindGVR from json file success: apps/v1, Resource=deployments
	//2022/09/02 14:30:13 FindGVR from yaml data success: apps/v1, Resource=deployments
	//2022/09/02 14:30:13 FindGVR from json data success: apps/v1, Resource=deployments
	//2022/09/02 14:30:13 FindGVR from map[string]interface{} success: apps/v1, Resource=deployments
	//2022/09/02 14:30:13 FindGVR from *unstructured.Unstructured success: apps/v1, Resource=deployments
	//2022/09/02 14:30:13 FindGVR from unstructured.Unstructured success: apps/v1, Resource=deployments
	//2022/09/02 14:30:13 FindGVR from runtime.Object success: apps/v1, Resource=deployments
	//2022/09/02 14:30:13 FindGVR from int failed: type must be string, []byte, map[string]interface{}, *unstructured.Unstructured, unstructured.Unstructured or runtime.Object
}
