package main

import (
	"io/ioutil"
	"log"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Get() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)
	handler.Apply(filename)

	deploy1, err := handler.GetByName(name)
	myerr("GetByName", err)

	deploy2, err := handler.Get(name)
	myerr("Get", err)

	deploy3, err := handler.GetFromFile(filename)
	myerr("GetFromFile", err)

	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	deploy4, err := handler.GetFromBytes(data)
	myerr("GetFromBytes", err)

	log.Println(deploy1.Name, deploy2.Name, deploy3.Name, deploy4.Name)

	// Output:

	//2022/07/04 21:51:05 GetByName success.
	//2022/07/04 21:51:05 Get success.
	//2022/07/04 21:51:05 GetFromFile success.
	//2022/07/04 21:51:05 GetFromBytes success.
	//2022/07/04 21:51:05 mydep mydep mydep mydep
}
