package main

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/forbearing/k8s/pod"
)

func Pod_Watch() {
	// New returns a handler used to multiples pod.
	handler, err := pod.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)
	handler.Apply(filename)

	addFunc := func(x interface{}) { log.Println("added pod.") }
	modifyFunc := func(x interface{}) { log.Println("modified pod.") }
	deleteFunc := func(x interface{}) { log.Println("deleted pod.") }

	// WatchByLabel watchs a set of pods by labels.
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

	// WatchByName watchs a pod by label.
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

	//2022/08/09 19:10:47 modified pod.
	//2022/08/09 19:10:47 modified pod.
	//2022/08/09 19:10:52 modified pod.
	//2022/08/09 19:11:04 modified pod.
	//2022/08/09 19:11:05 modified pod.
	//2022/08/09 19:11:05 deleted pod.
	//2022/08/09 19:11:07 added pod.
	//2022/08/09 19:11:07 modified pod.
	//2022/08/09 19:11:07 modified pod.
	//2022/08/09 19:11:12 modified pod.
	//2022/08/09 19:11:24 modified pod.
	//2022/08/09 19:11:24 modified pod.
	//2022/08/09 19:11:25 modified pod.
	//2022/08/09 19:11:25 modified pod.
	//2022/08/09 19:11:25 deleted pod.
	//2022/08/09 19:11:25 deleted pod.
	//2022/08/09 19:11:27 added pod.
	//2022/08/09 19:11:27 added pod.
	//2022/08/09 19:11:27 modified pod.
	//2022/08/09 19:11:27 modified pod.
	//2022/08/09 19:11:27 modified pod.
	//2022/08/09 19:11:27 modified pod.
	//2022/08/09 19:11:27 modified pod.
	//2022/08/09 19:11:27 modified pod.
	//2022/08/09 19:11:43 modified pod.
	//2022/08/09 19:11:43 modified pod.
	//2022/08/09 19:11:44 modified pod.
	//2022/08/09 19:11:44 modified pod.
	//2022/08/09 19:11:44 deleted pod.
	//2022/08/09 19:11:44 deleted pod.
}
