package main

import (
	"testing"

	"github.com/forbearing/k8s/deployment"
)

func BenchmarkHandler(b *testing.B) {
	_, err := deployment.New(ctx, kubeconfig, namespace)
	if err != nil {
		b.Error(err)
	}
}
