package main

import (
	"fmt"
	"log"
	"time"

	"github.com/forbearing/k8s/job"
)

func Job_Apply() {
	handler := job.NewOrDie(ctx, "", namespace)

	if _, err := handler.Apply(filename); err != nil {
		log.Fatal(err)
	}
	jobObj, err := handler.Get(name)
	if err != nil {
		log.Fatal(err)
	}

	annotations := jobObj.Spec.Template.Annotations
	if annotations == nil {
		annotations = make(map[string]string)
	}
	fmt.Println()
	fmt.Println()

	annotations["updatedTimestamp"] = time.Now().Format(time.RFC3339)
	if _, err := handler.Apply(jobObj); err != nil {
		log.Fatal(err)
	}

	if err := handler.Delete(name); err != nil {
		log.Fatalf("delete job error: %v", err)
	}
}
