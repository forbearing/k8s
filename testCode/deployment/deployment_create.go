package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
)

var rawName = "mydep-raw"
var rawData = map[string]interface{}{
	"apiVersion": "apps/v1",
	"kind":       "Deployment",
	"metadata": map[string]interface{}{
		"name": rawName,
		"labels": map[string]interface{}{
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

func Deployment_Create() {
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer create(handler)

	// CreateFromRaw
	_, err = handler.CreateFromRaw(rawData)
	myerr("CreateFromRaw", err)
	handler.Delete(name)

	// CreateFromFile
	_, err = handler.CreateFromFile(filename)
	myerr("CreateFromFile", err)
	handler.Delete(name)

	// CreateFromBytes
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	_, err = handler.CreateFromBytes(data)
	myerr("CreateFromBytes", err)
	handler.Delete(name)

	// Create
	_, err = handler.Create(filename)
	myerr("Create", err)

}
