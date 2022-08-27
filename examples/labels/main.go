package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/util/labels"
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
	// has label?
	log.Println(labels.Has(deploy, "app"))         // true
	log.Println(labels.Has(deploy, "app=mydep"))   // true
	log.Println(labels.Has(deploy, "app=novalue")) // false
	log.Println(labels.Has(deploy, "nolabel"))     // false

	// get all labels
	log.Println(labels.GetAll(deploy)) // map[app:mydep type:deployment]
	// get a label
	log.Println(labels.Get(deploy, "app"))         // mydep
	log.Println(labels.Get(deploy, "app=mydep"))   // mydep
	log.Println(labels.Get(deploy, "app=novalue")) // mydep

	// set label
	fmt.Println()
	log.Println("before labels: ", labels.GetAll(deploy))
	labels.Set(deploy, "label1=value1", "label2=value2")
	log.Println("after labels : ", labels.GetAll(deploy))
	labels.Set(deploy, "label1=value1", "label2=value2")
	log.Println("after labels : ", labels.GetAll(deploy))
	// Output:
	//2022/08/27 16:06:54 before labels:  map[app:mydep type:deployment]
	//2022/08/27 16:06:54 after labels :  map[app:mydep label1:value1 label2:value2 type:deployment]
	//2022/08/27 16:06:54 after labels :  map[app:mydep label1:value1 label2:value2 type:deployment]

	// remove label
	fmt.Println()
	log.Println("before labels: ", labels.GetAll(deploy))
	labels.Remove(deploy, "label1=value1", "label2=value2")
	log.Println("after labels : ", labels.GetAll(deploy))
	labels.Remove(deploy, "label1=value1", "label2=value2")
	log.Println("after labels : ", labels.GetAll(deploy))
	// Output
	//2022/08/27 16:07:42 before labels:  map[app:mydep label1:value1 label2:value2 type:deployment]
	//2022/08/27 16:07:42 after labels :  map[app:mydep type:deployment]
	//2022/08/27 16:07:42 after labels :  map[app:mydep type:deployment]

	// remove all labels
	fmt.Println()
	log.Println("before labels: ", labels.GetAll(deploy))
	labels.RemoveAll(deploy)
	log.Println("after labels: ", labels.GetAll(deploy))
	// Output:
	//2022/08/27 16:08:56 before labels:  map[app:mydep type:deployment]
	//2022/08/27 16:08:56 after labels:  map[]
}

// cleanup will delete or prune created deployments.
func cleanup(handler *deployment.Handler) {
	handler.Delete(name)
}
