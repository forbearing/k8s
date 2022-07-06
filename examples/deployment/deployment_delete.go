package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Delete() {
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// DeleteByName delete a deployment by name.
	handler.Apply(filename)
	myerr("DeleteByName", handler.DeleteByName(name))

	// Delete delete a deployment by name, it's alias to "DeleteByName".
	handler.Apply(filename)
	myerr("Delete", handler.Delete(name))

	// DeleteFromFile delete a deployment from yaml file.
	handler.Apply(filename)
	myerr("DeleteFromFile", handler.DeleteFromFile(filename))

	// DeleteFromBytes delete a deployment from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	handler.Apply(filename)
	myerr("DeleteFromBytes", handler.DeleteFromBytes(data))

	// Output:

	//2022/07/04 21:43:08 DeleteByName success.
	//2022/07/04 21:43:08 Delete success.
	//2022/07/04 21:43:08 DeleteFromFile success.
	//2022/07/04 21:43:08 DeleteFromBytes success.
}
