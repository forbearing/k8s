package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/util/annotations"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	filename    = "../../testdata/examples/deployment.yaml"
	name        = "mydep"
)

func main() {
	handler := deployment.NewOrDie(ctx, "", namespace)
	defer cleanup(handler)

	deploy, err := handler.Apply(filename)
	if err != nil {
		log.Fatal(err)
	}
	// has annotation?
	log.Println(annotations.Has(deploy, "region"))      // true
	log.Println(annotations.Has(deploy, "region=east")) // true
	log.Println(annotations.Has(deploy, "region=west")) // false
	log.Println(annotations.Has(deploy, "zone"))        // false

	// get all annotations
	log.Println(annotations.GetAll(deploy)) // map[region:east]
	// get a annotation
	log.Println(annotations.Get(deploy, "region"))      // east
	log.Println(annotations.Get(deploy, "region=east")) // east
	log.Println(annotations.Get(deploy, "region=west")) // east

	// set annotation
	fmt.Println()
	log.Println("before annotations: ", annotations.GetAll(deploy))
	annotations.Set(deploy, "zone=us", "master=true")
	log.Println("after annotations : ", annotations.GetAll(deploy))
	annotations.Set(deploy, "zone=us", "master=true")
	log.Println("after annotations : ", annotations.GetAll(deploy))
	// Output:
	//2022/08/27 17:17:34 before annotations:  map[region:east]
	//2022/08/27 17:17:34 after annotations :  map[master:true region:east zone:us]
	//2022/08/27 17:17:34 after annotations :  map[master:true region:east zone:us]

	// remove annotation
	fmt.Println()
	log.Println("before annotations: ", annotations.GetAll(deploy))
	annotations.Remove(deploy, "zone=us", "master=true")
	log.Println("after annotations : ", annotations.GetAll(deploy))
	annotations.Remove(deploy, "zone=us", "master=true")
	log.Println("after annotations : ", annotations.GetAll(deploy))
	// Output
	//2022/08/27 17:18:27 before annotations:  map[master:true region:east zone:us]
	//2022/08/27 17:18:27 after annotations :  map[region:east]
	//2022/08/27 17:18:27 after annotations :  map[region:east]

	// remove all annotations
	fmt.Println()
	log.Println("before annotations: ", annotations.GetAll(deploy))
	annotations.RemoveAll(deploy)
	log.Println("after annotations: ", annotations.GetAll(deploy))
	// Output:
	//2022/08/27 17:19:00 before annotations:  map[region:east]
	//2022/08/27 17:19:00 after annotations:  map[]
}

// cleanup will delete or prune created deployments.
func cleanup(handler *deployment.Handler) {
	handler.Delete(name)
}
