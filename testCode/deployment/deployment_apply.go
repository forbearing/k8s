package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/deployment"
)

// ApplyXXX method will creates deployment if specified deployment not exist and
// it will updates deployment if already exist.
func Deployment_Apply() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	handler.Delete(name)
	handler.Delete(rawName)

	// ApplyFromRaw apply a deployment from map[string]interface.
	// it wil updates the deployment, if already exist, or creates it.
	handler.CreateFromRaw(rawData)
	_, err = handler.ApplyFromRaw(rawData)
	myerr("ApplyFromRaw", err)

	handler.Delete(rawName)
	_, err = handler.ApplyFromRaw(rawData)
	myerr("ApplyFromRaw", err)

	// ApplyFromFile apply a deployment from file.
	// it wil updates the deployment, if already exist, or creates it.
	handler.CreateFromFile(update1File)
	_, err = handler.ApplyFromFile(update1File)
	myerr("ApplyFromFile", err)
	handler.DeleteFromFile(update1File)
	_, err = handler.ApplyFromFile(update1File)
	myerr("ApplyFromFile", err)

	// ApplyFromBytes apply a deployment from bytes.
	// it wil updates the deployment, if already exist, or creates it.
	var data []byte
	if data, err = ioutil.ReadFile(update2File); err != nil {
		panic(err)
	}
	handler.CreateFromFile(update2File)
	_, err = handler.ApplyFromBytes(data)
	myerr("ApplyFromBytes", err)
	handler.DeleteFromFile(update2File)
	_, err = handler.ApplyFromBytes(data)
	myerr("ApplyFromBytes", err)

	// Apply apply a deployment from file, it's alias to "ApplyFromFile".
	// it wil updates the deployment, if already exist, or creates it.
	handler.CreateFromFile(update3File)
	_, err = handler.Apply(update3File)
	myerr("Apply", err)
	handler.DeleteFromFile(update3File)
	myerr("Apply", err)

	// Output:

	//2022/07/04 21:43:05 ApplyFromRaw success.
	//2022/07/04 21:43:05 ApplyFromRaw success.
	//2022/07/04 21:43:05 ApplyFromFile success.
	//2022/07/04 21:43:05 ApplyFromFile success.
	//2022/07/04 21:43:06 ApplyFromBytes success.
	//2022/07/04 21:43:06 ApplyFromBytes success.
	//2022/07/04 21:43:07 Apply success.
	//2022/07/04 21:43:07 Apply success.

}
