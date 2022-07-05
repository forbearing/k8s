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
	filename    = "./testData/examples/all.yaml"
)

func TestK8S(t *testing.T) {
	defer cancel()

	err := ApplyF(ctx, kubeconfig, filename)
	myerr(t, "ApplyF", err)

	err = DeleteF(ctx, kubeconfig, filename)
	myerr(t, "DeleteF", err)
}

func myerr(t *testing.T, name string, err error) {
	if err != nil {
		t.Errorf("%s failed: %v", name, err)
	} else {
		t.Logf("%s success.", name)
	}
}