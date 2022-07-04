package daemonset

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../testData/examples/daemonset.yaml"
	name        = "myds"
	label       = "type=daemonset"
	rawName     = "myds-raw"
)

func TestDaemonset(t *testing.T) {
	defer cancel()

	t.Run("Create Daemonset", testCreateDaemonset)
	t.Run("Update Daemonset", testUpdateDaemonset)
	t.Run("Apply Daemonset", testApplyDaemonset)
	t.Run("Delete Daemonset", testDeleteDaemonset)
	t.Run("Get Daemonset", testGetDaemonset)
	t.Run("List Daemonset", testListDaemonset)
	t.Run("Daemonset Tools", testDaemonsetTools)
}

func testCreateDaemonset(t *testing.T) {}
func testUpdateDaemonset(t *testing.T) {}
func testApplyDaemonset(t *testing.T)  {}
func testDeleteDaemonset(t *testing.T) {}
func testGetDaemonset(t *testing.T)    {}
func testListDaemonset(t *testing.T)   {}
func testWatchDaemonset(t *testing.T)  {}
func testDaemonsetTools(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	// test IsReady, WaitReady
	_, err = handler.Apply(filename)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s is ready: %v", name, handler.IsReady(name))
	handler.WaitReady(name)
	t.Logf("%s is ready: %v", name, handler.IsReady(name))
	handler.Delete(name)

	// test GetPods
	name = "nginx-ds"
	podList, err := handler.GetPods(name)
	myerr(t, "GetPods", err)
	outputPods(t, podList)

	// test GetPVC
	pvcList, err := handler.GetPVC(name)
	myerr(t, "GetPVC", err)
	t.Log(pvcList)

	// test GetPV
	pvList, err := handler.GetPV(name)
	myerr(t, "GetPV", err)
	t.Log(pvList)

	// test GetAge
	age, err := handler.GetAge(name)
	myerr(t, "GetAge", err)
	t.Log(age)
}

func myerr(t *testing.T, name string, err error) {
	if err != nil {
		t.Errorf("%s failed: %v", name, err)
	} else {
		t.Logf("%s success.", name)
	}
}
func outputPods(t *testing.T, podList []corev1.Pod) {
	var pl []string
	for _, p := range podList {
		pl = append(pl, p.Name)
	}
	t.Log(pl)
}
