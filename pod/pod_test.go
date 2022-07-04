package pod

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../testData/examples/pod.yaml"
	name        = "mypod"
	label       = "type=pod"
	rawName     = "mypod-raw"
)

func TestPod(t *testing.T) {
	defer cancel()

	t.Run("Create Pod", testCreatePod)
	t.Run("Update Pod", testUpdatePod)
	t.Run("Apply Pod", testApplyPod)
	t.Run("Delete Pod", testDeletePod)
	t.Run("Get Pod", testGetPod)
	t.Run("List Pod", testListPod)
	t.Run("Pod Tools", testPodTools)
}

func testCreatePod(t *testing.T) {}
func testUpdatePod(t *testing.T) {}
func testApplyPod(t *testing.T)  {}
func testDeletePod(t *testing.T) {}
func testGetPod(t *testing.T)    {}
func testListPod(t *testing.T)   {}
func testPodTools(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	handler.Apply(filename)

	// test IsReady, WaitReady
	t.Logf("%s is ready: %t", name, handler.IsReady(name))
	handler.WaitReady(name)
	t.Logf("%s is ready: %v", name, handler.IsReady(name))

	// test GetUID, GetIP
	uid, err := handler.GetUID(name)
	myerr(t, "GetUID", err)
	t.Log(uid)
	ip, err := handler.GetIP(name)
	myerr(t, "GetIP", err)
	t.Log(ip)

	// test GetNodeName, GetNodeIP
	nodeName, err := handler.GetNodeName(name)
	myerr(t, "GetNodeName", err)
	t.Log(nodeName)
	nodeIP, err := handler.GetNodeIP(name)
	myerr(t, "GetNodeIP", err)
	t.Log(nodeIP)

	// test GetAge, GetStatus, GetQosClass
	age, err := handler.GetAge(name)
	myerr(t, "GetAge", err)
	t.Log(age)
	status, err := handler.GetStatus(name)
	myerr(t, "GetStatus", err)
	t.Log(status)
	qos, err := handler.GetQosClass(name)
	myerr(t, "GetQosClass", err)
	t.Log(qos)

	// test GetContainers, GetInitContainers, GetReadyContainers
	cl, err := handler.GetContainers(name)
	myerr(t, "GetContainers", err)
	t.Log(cl)
	icl, err := handler.GetInitContainers(name)
	myerr(t, "GetInitContainers", err)
	t.Log(icl)
	rcl, err := handler.GetReadyContainers(name)
	myerr(t, "GetReadyContainers", err)
	t.Log(rcl)

	// test GetPVC, GetPV
	name2 := "nginx-pod"
	pvcList, err := handler.GetPVC(name2)
	myerr(t, "GetPVC", err)
	t.Log(pvcList)
	pvList, err := handler.GetPV(name2)
	myerr(t, "GetPV", err)
	t.Log(pvList)

	// test GetController
	pc1, err := handler.GetController(name)
	// will test fail
	myerr(t, "GetController", err)
	t.Log(pc1)
	// will test success.
	// you should execute "kubectl apply -f testData/nginx/nginx-sts.yaml",
	// before get controller of "nginx-sts-0".
	pc2, err := handler.GetController("nginx-sts-0")
	myerr(t, "GetController", err)
	t.Log(pc2)

	// test Execute
	command := []string{
		"/bin/sh",
		"-c",
		"cat /etc/os-release",
	}
	err = handler.Execute(name, "", command)
	myerr(t, "Execute", err)

	//handler.DeleteFromFile(filename)
}
func myerr(t *testing.T, name string, err error) {
	if err != nil {
		t.Errorf("%s failed: %v", name, err)
	} else {
		t.Logf("%s success.", name)
	}
}
