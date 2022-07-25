package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
)

func Deployment_Watch() {
	// New returns a handler used to multiples deployment.
	handler, err := deployment.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)
	handler.Apply(filename)

	addFunc := func(x interface{}) { log.Println("added deployment.") }
	modifyFunc := func(x interface{}) { log.Println("modified deployment.") }
	deleteFunc := func(x interface{}) { log.Println("deleted deployment.") }

	// WatchByLabel watchs a set of deployments by label.
	{
		ctx, cancel := context.WithCancel(ctx)

		go func(ctx context.Context) {
			handler.WatchByLabel(label, addFunc, modifyFunc, deleteFunc, nil)
		}(ctx)
		go func(ctx context.Context) {
			for {
				handler.Apply(filename)
				time.Sleep(time.Second * 5)
				handler.Delete(name)
			}
		}(ctx)

		timer := time.NewTimer(time.Second * 30)
		<-timer.C
		cancel()
	}

	// WatchByName watchs a deployment by label.
	// Watch simply calls WatchByName.
	ioutil.ReadFile(filename)
	{
		ctx, cancel := context.WithCancel(ctx)

		go func(ctx context.Context) {
			handler.WatchByName(name, addFunc, modifyFunc, deleteFunc, nil)
		}(ctx)
		go func(ctx context.Context) {
			for {
				handler.Apply(filename)
				time.Sleep(time.Second * 5)
				handler.Delete(name)
			}
		}(ctx)

		timer := time.NewTimer(time.Second * 30)
		<-timer.C
		cancel()
	}

	// Output:

	//2022/07/04 21:59:58 modified deployment.
	//2022/07/04 21:59:58 modified deployment.
	//2022/07/04 21:59:58 modified deployment.
	//2022/07/04 21:59:58 modified deployment.
	//2022/07/04 21:59:58 modified deployment.
	//2022/07/04 22:00:03 deleted deployment.
	//2022/07/04 22:00:03 added deployment.
	//2022/07/04 22:00:03 modified deployment.
	//2022/07/04 22:00:03 modified deployment.
	//2022/07/04 22:00:03 modified deployment.
	//2022/07/04 22:00:08 deleted deployment.
	//2022/07/04 22:00:08 added deployment.
	//2022/07/04 22:00:08 modified deployment.
	//2022/07/04 22:00:08 modified deployment.
	//2022/07/04 22:00:08 modified deployment.
	//2022/07/04 22:00:08 modified deployment.
	//2022/07/04 22:00:13 deleted deployment.
	//2022/07/04 22:00:13 added deployment.
	//2022/07/04 22:00:13 modified deployment.
	//2022/07/04 22:00:13 modified deployment.
	//2022/07/04 22:00:13 modified deployment.
	//2022/07/04 22:00:18 deleted deployment.
	//2022/07/04 22:00:18 added deployment.
	//2022/07/04 22:00:18 modified deployment.
	//2022/07/04 22:00:18 modified deployment.
	//2022/07/04 22:00:18 modified deployment.
	//2022/07/04 22:00:23 deleted deployment.
	//2022/07/04 22:00:23 added deployment.
	//2022/07/04 22:00:23 modified deployment.
	//2022/07/04 22:00:23 modified deployment.
	//2022/07/04 22:00:23 modified deployment.
	//2022/07/04 22:00:28 deleted deployment.
}
