package statefulset

import (
	"context"
	"io/ioutil"
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
	filename    = "../testData/examples/statefulset.yaml"
	name        = "mysts"
	label       = "type=statefulset"
	rawName     = "mysts-raw"
)
var rawData = map[string]interface{}{
	"apiVersion": "apps/v1",
	"kind":       "Statefulset",
	"metadata": map[string]interface{}{
		"name": rawName,
		"labels": map[string]interface{}{
			"app":  rawName,
			"type": "statefulset",
		},
	},
	"spec": map[string]interface{}{
		"replicas":    3,
		"serviceName": "sts-headless",
		"selector": map[string]interface{}{
			"matchLabels": map[string]interface{}{
				"app":  rawName,
				"type": "statefulset",
			},
		},
		"template": map[string]interface{}{
			"metadata": map[string]interface{}{
				"labels": map[string]interface{}{
					"app":  rawName,
					"type": "statefulset",
				},
			},
			"spec": map[string]interface{}{
				"containers": []map[string]interface{}{
					{
						"name":  "nginx",
						"image": "nginx",
					},
				},
			},
		},
	},
}

func TestStatefulSet(t *testing.T) {
	defer cancel()

	//t.Run("Create Statefulset", testCreateStatefulset)
	//t.Run("Update Statefulset", testUpdateStatefulset)
	//t.Run("Apply Statefulset", testApplyStatefulset)
	//t.Run("Delete Statefulset", testDeleteStatefulset)
	//t.Run("Get Statefulset", testGetStatefulSet)
	//t.Run("List Statefulset", testListStatefulset)
	//t.Run("Watch Statefulset", testWatchStatefulset)
	t.Run("Statefulset Tools", testStatefulsetTools)
}

func testCreateStatefulset(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	handler.Delete(name)

	sts, err := handler.Create(filename)
	myerr(t, "Create", err)
	handler.Delete(sts.Name)

	sts, err = handler.CreateFromFile(filename)
	myerr(t, "CreateFromFile", err)
	handler.Delete(sts.Name)

	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		t.Fatal(err)
	}
	sts, err = handler.CreateFromBytes(data)
	myerr(t, "CreateFromBytes", err)
	handler.Delete(sts.Name)

	sts, err = handler.CreateFromRaw(rawData)
	myerr(t, "CreateFromRaw", err)
	handler.Delete(sts.Name)
}
func testUpdateStatefulset(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	handler.Create(filename)

	filename := "../testData/examples/statefulset-update1.yaml"
	sts, err := handler.Update(filename)
	myerr(t, "Update", err)

	filename = "../testData/examples/statefulset-update2.yaml"
	sts, err = handler.UpdateFromFile(filename)
	myerr(t, "UpdateFromFile", err)

	filename = "../testData/examples/statefulset-update3.yaml"
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		t.Fatal(err)
	}

	sts, err = handler.UpdateFromBytes(data)
	myerr(t, "UpdateFromBytes", err)
	handler.Delete(sts.Name)

	handler.CreateFromRaw(rawData)
	sts, err = handler.UpdateFromRaw(rawData)
	myerr(t, "UpdateFromRaw", err)
	handler.Delete(sts.Name)

}
func testApplyStatefulset(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	handler.Create(filename)

	filename := "../testData/examples/statefulset-update1.yaml"
	sts, err := handler.Apply(filename)
	myerr(t, "Apply", err)
	handler.Delete(sts.Name)
	sts, err = handler.Apply(filename)
	myerr(t, "Apply", err)
	handler.Delete(sts.Name)

	filename = "../testData/examples/statefulset-update2.yaml"
	sts, err = handler.ApplyFromFile(filename)
	myerr(t, "ApplyFromFile", err)
	handler.Delete(sts.Name)
	sts, err = handler.ApplyFromFile(filename)
	myerr(t, "ApplyFromFile", err)
	handler.Delete(sts.Name)

	filename = "../testData/examples/statefulset-update3.yaml"
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		t.Fatal(err)
	}
	sts, err = handler.ApplyFromBytes(data)
	myerr(t, "ApplyFromBytes", err)
	handler.Delete(sts.Name)
	sts, err = handler.ApplyFromBytes(data)
	myerr(t, "ApplyFromBytes", err)
	handler.Delete(sts.Name)

	sts, err = handler.ApplyFromRaw(rawData)
	myerr(t, "ApplyFromRaw", err)
	handler.Delete(sts.Name)
	sts, err = handler.ApplyFromRaw(rawData)
	myerr(t, "ApplyFromRaw", err)
	handler.Delete(sts.Name)
}
func testDeleteStatefulset(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	handler.Apply(filename)
	myerr(t, "Delete", handler.Delete(name))

	handler.Apply(filename)
	myerr(t, "DeleteByName", handler.DeleteByName(name))

	handler.Apply(filename)
	myerr(t, "DeleteFromFile", handler.DeleteFromFile(filename))

	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		t.Fatal(err)
	}
	handler.Apply(filename)
	myerr(t, "DeleteFromBytes", handler.DeleteFromBytes(data))
}
func testGetStatefulSet(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	handler.Apply(filename)

	sts1, err := handler.Get(name)
	myerr(t, "Get", err)

	sts2, err := handler.GetByName(name)
	myerr(t, "GetByName", err)

	sts3, err := handler.GetFromFile(filename)
	myerr(t, "GetFromFile", err)

	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		t.Fatal(err)
	}
	sts4, err := handler.GetFromBytes(data)
	myerr(t, "GetFromBytes", err)
	t.Log(sts1.Name, sts2.Name, sts3.Name, sts4.Name)

	handler.Delete(name)
}
func testListStatefulset(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	handler.Apply("../testData/examples/statefulset.yaml")
	handler.Apply("../testData/examples/statefulset-2.yaml")
}
func testWatchStatefulset(t *testing.T) {}

func testStatefulsetTools(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	handler.Apply(filename)

	// test IsReady, WaitReady
	t.Logf("statefulset is ready: %t", handler.IsReady(name))
	handler.WaitReady(name)
	t.Logf("statefulset is ready: %t", handler.IsReady(name))

	// test GetPods
	podList, err := handler.GetPods(name)
	myerr(t, "GetPods", err)
	outputPods(t, podList)
	handler.DeleteFromFile(filename)

	// test GetPVC
	name := "nginx-sts"
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
	for _, pod := range podList {
		pl = append(pl, pod.Name)
	}
	t.Log(pl)
}
