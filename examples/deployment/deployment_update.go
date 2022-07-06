package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Update() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)
	handler.Apply(filename)
	handler.ApplyFromRaw(rawData)

	// UpdateFromRaw updates deployment from map[string]interface.
	_, err = handler.UpdateFromRaw(rawData)
	myerr("UpdateFromRaw", err)

	// UpdateFromFile updates deployment from file.
	_, err = handler.UpdateFromFile(update1File)
	myerr("UpdateFromFile", err)

	// UpdateFromBytes updates deployment from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(update2File); err != nil {
		panic(err)
	}
	_, err = handler.UpdateFromBytes(data)
	myerr("UpdateFromBytes", err)

	// Update updates deployment from file, it's alias to "UpdateFromFile".
	_, err = handler.Update(update3File)
	myerr("Update", err)

	// Output:

	//2022/07/04 21:43:05 UpdateFromRaw success.
	//2022/07/04 21:43:05 UpdateFromFile success.
	//2022/07/04 21:43:05 UpdateFromBytes success.
	//2022/07/04 21:43:05 Update success.
}
