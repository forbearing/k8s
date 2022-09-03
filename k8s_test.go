package k8s

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "./testdata/examples/all.yaml"
)

func TestK8S(t *testing.T) {
	defer cancel()

	namespace := "test"

	err := ApplyF(ctx, kubeconfig, filename, namespace)
	checkErr(t, "ApplyF", err)

	err = DeleteF(ctx, kubeconfig, filename, namespace)
	checkErr(t, "DeleteF", err)
}

func checkErr(t *testing.T, name string, err error) {
	if err != nil {
		t.Errorf("%s failed: %v", name, err)
	} else {
		t.Logf("%s success.", name)
	}
}
