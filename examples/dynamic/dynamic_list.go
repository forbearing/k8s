package main

import (
	"context"
	"fmt"
	"log"

	"github.com/forbearing/k8s/deployment"
	"github.com/forbearing/k8s/dynamic"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Dynamic_List() {
	filename := "../../testdata/examples/deployment.yaml"
	handler := dynamic.NewOrDie(context.TODO(), "", namespace)
	defer cleanup(handler)

	if _, err := handler.Apply(deployUnstructData); err != nil {
		log.Fatal(err)
	}
	if _, err := handler.Apply(filename); err != nil {
		log.Fatal(err)
	}

	gvk := deployment.GVK()
	label := "type=deployment"
	field := fmt.Sprintf("metadata.namespace=%s", namespace)
	_ = label
	_ = field

	ul, err := handler.WithGVK(gvk).List()
	checkErr("List()", output(ul), err)
	ul2, err := handler.WithGVK(gvk).ListAll()
	checkErr("ListAll()", output(ul2), err)
	ul3, err := handler.WithGVK(gvk).ListByLabel(label)
	checkErr("ListByLabel()", output(ul3), err)
	ul4, err := handler.WithGVK(gvk).ListByField(field)
	checkErr("ListByField()", output(ul4), err)
	ul5, err := handler.WithGVK(gvk).ListByNamespace(namespace)
	checkErr("ListByNamespace()", output(ul5), err)

	// Output:

	//2022/09/07 18:15:19 List() success: [nginx coredns local-path-provisioner mydep mydep-unstruct]
	//2022/09/07 18:15:19 ListAll() success: [nginx coredns local-path-provisioner mydep mydep-unstruct]
	//2022/09/07 18:15:19 ListByLabel() success: [mydep mydep-unstruct]
	//2022/09/07 18:15:19 ListByField() success: [mydep mydep-unstruct]
	//2022/09/07 18:15:19 ListByNamespace() success: [mydep mydep-unstruct]
}

func output(unstructList []*unstructured.Unstructured) []string {
	if unstructList == nil {
		return nil
	}
	var ul []string
	for _, u := range unstructList {
		ul = append(ul, u.GetName())
	}
	return ul
}
