package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Create() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// CreateFromRaw creates a deploymnt from map[string]interface.
	_, err = handler.CreateFromRaw(rawData)
	myerr("CreateFromRaw", err)
	handler.Delete(name)

	// CreateFromFile creates a deploymnt from file.
	_, err = handler.CreateFromFile(filename)
	myerr("CreateFromFile", err)
	handler.Delete(name)

	// CreateFromBytes creates a deploymnt from bytes.
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	_, err = handler.CreateFromBytes(data)
	myerr("CreateFromBytes", err)
	handler.Delete(name)

	// Create creates a deploymnt from file, it's alias to "CreateFromFile".
	_, err = handler.Create(filename)
	myerr("Create", err)

	// Output:

	//2022/07/04 21:43:04 CreateFromRaw success.
	//2022/07/04 21:43:04 CreateFromFile success.
	//2022/07/04 21:43:04 CreateFromBytes success.
	//2022/07/04 21:43:04 Create success.
}
