package deployment

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	filename    = "../testData/examples/deployment.yaml"
	label       = "app=mydep"
	name1       = "mydep"
	name2       = "mydep-raw"
)

var rawData1 = map[string]interface{}{
	"apiVersion": "apps/v1",
	"kind":       "Deployment",
	"metadata": map[string]interface{}{
		"name": name1,
		"labels": map[string]interface{}{
			"type": "deployment",
		},
	},
	"spec": map[string]interface{}{
		// replicas type is int32, not string.
		"replicas": 2,
		"selector": map[string]interface{}{
			"matchLabels": map[string]interface{}{
				"app":  name1,
				"type": "deployment",
			},
		},
		"template": map[string]interface{}{
			"metadata": map[string]interface{}{
				"labels": map[string]interface{}{
					"app":  name1,
					"type": "deployment",
				},
			},
			"spec": map[string]interface{}{
				"containers": []map[string]interface{}{
					{
						"name":  "nginx",
						"image": "nginx",
						"resources": map[string]interface{}{
							"limits": map[string]interface{}{
								"cpu": "100m",
							},
						},
					},
				},
			},
		},
	},
}
var rawData2 = map[string]interface{}{
	"apiVersion": "apps/v1",
	"kind":       "Deployment",
	"metadata": map[string]interface{}{
		"name": name2,
		"labels": map[string]interface{}{
			"type": "deployment",
		},
	},
	"spec": map[string]interface{}{
		// replicas type is int32, not string.
		"replicas": 1,
		"selector": map[string]interface{}{
			"matchLabels": map[string]interface{}{
				"app":  name2,
				"type": "deployment",
			},
		},
		"template": map[string]interface{}{
			"metadata": map[string]interface{}{
				"labels": map[string]interface{}{
					"app":  name2,
					"type": "deployment",
				},
			},
			"spec": map[string]interface{}{
				"containers": []map[string]interface{}{
					{
						"name":  "nginx",
						"image": "nginx",
						"resources": map[string]interface{}{
							"limits": map[string]interface{}{
								"cpu": "100m",
							},
						},
					},
				},
			},
		},
	},
}

func TestDeployment(t *testing.T) {
	defer cancel()

	//t.Run("Create Deployment", testCreateDeployment)
	//t.Run("Update Deployment", testUpdateDeployment)
	//t.Run("Apply Deployment", testApplyDeployment)
	//t.Run("Delete Deployment", testDeleteDeployment)
	//t.Run("Get Deployment", testGetDeployment)
	//t.Run("List Deployment", testListDeployment)
	//t.Run("Watch Deployment", testWatchDeployment)
	t.Run("Deployment Tools", testDeploymentTools)
}

func testCreateDeployment(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	handler.DeleteFromFile(filename)
	_, err = handler.Create(filename)
	myerr(t, "Create", err)

	handler.DeleteFromFile(filename)
	_, err = handler.CreateFromFile(filename)
	myerr(t, "CreateFromFile", err)

	handler.DeleteFromFile(filename)
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		t.Fatal(err)
	}
	_, err = handler.CreateFromBytes(data)
	myerr(t, "CreateFromBytes", err)
	handler.DeleteFromFile(filename)

	handler.Delete(name2)
	_, err = handler.CreateFromRaw(rawData2)
	myerr(t, "CreateFromRaw", err)
	handler.Delete(name2)
}

func testUpdateDeployment(t *testing.T) {
	var deploy *appsv1.Deployment
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	deploy, _ = handler.Apply(filename)

	filename = "../testData/examples/deployment-update1.yaml"
	_, err = handler.Update(filename)
	myerr(t, "Update", err)

	filename = "../testData/examples/deployment-update2.yaml"
	_, err = handler.UpdateFromFile(filename)
	myerr(t, "UpdateFromFile", err)

	filename = "../testData/examples/deployment-update3.yaml"
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		t.Fatal(err)
	}
	_, err = handler.UpdateFromBytes(data)
	myerr(t, "UpdateFromBytes", err)

	_, err = handler.UpdateFromRaw(rawData1)
	myerr(t, "UpdateFromRaw", err)
	handler.Delete(deploy.Name)
}

func testApplyDeployment(t *testing.T) {
	var deploy *appsv1.Deployment
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	deploy, _ = handler.Apply(filename)
	filename = "../testData/examples/deployment-update1.yaml"
	_, err = handler.Apply(filename)
	myerr(t, "Apply", err)
	handler.Delete(deploy.Name)
	_, err = handler.Apply(filename)
	myerr(t, "Apply", err)

	filename = "../testData/examples/deployment-update2.yaml"
	_, err = handler.ApplyFromFile(filename)
	myerr(t, "ApplyFromFile", err)
	handler.Delete(deploy.Name)
	_, err = handler.ApplyFromFile(filename)
	myerr(t, "ApplyFromFile", err)

	filename = "../testData/examples/deployment-update3.yaml"
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		t.Fatal(err)
	}
	_, err = handler.ApplyFromBytes(data)
	myerr(t, "ApplyFromBytes", err)
	handler.Delete(deploy.Name)
	_, err = handler.ApplyFromBytes(data)
	myerr(t, "ApplyFromBytes", err)

	deploy, err = handler.ApplyFromRaw(rawData2)
	myerr(t, "ApplyFromRaw", err)
	fmt.Println(deploy.Name)
	deploy, err = handler.ApplyFromRaw(rawData2)
	myerr(t, "ApplyFromRaw", err)
	handler.Delete(deploy.Name)

}

func testDeleteDeployment(t *testing.T) {
	var deploy *appsv1.Deployment
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	deploy, err = handler.Apply(filename)
	myerr(t, "Delete", handler.Delete(deploy.Name))

	handler.Apply(filename)
	myerr(t, "DeleteByName", handler.DeleteByName(deploy.Name))

	handler.Apply(filename)
	myerr(t, "DeleteFromFile", handler.DeleteFromFile(filename))

	handler.Apply(filename)
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		t.Fatal(err)
	}
	myerr(t, "DeleteFromBytes", handler.DeleteFromBytes(data))
}

func testGetDeployment(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	_, err = handler.Apply(filename)
	if err != nil {
		t.Fatal(err)
	}

	deploy1, err := handler.Get(name1)
	myerr(t, "Get", err)

	deploy2, err := handler.GetByName(name1)
	myerr(t, "GetByName", err)

	deploy3, err := handler.GetFromFile(filename)
	myerr(t, "GetFromFile", err)

	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		t.Fatal(err)
	}
	deploy4, err := handler.GetFromBytes(data)
	myerr(t, "GetFromBytes", err)

	t.Log(deploy1.Name, deploy2.Name, deploy3.Name, deploy4.Name)
}

func testListDeployment(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	filename2 := "../testData/examples/deployment-2.yaml"
	_, err = handler.Apply(filename)
	if err != nil {
		t.Fatal(err)
	}
	_, err = handler.Apply(filename2)
	if err != nil {
		t.Fatal(err)
	}

	deployList1, err := handler.List(label)
	myerr(t, "List", err)
	outputDeploy(t, deployList1)

	deployList2, err := handler.ListByLabel(label)
	myerr(t, "ListByLabel", err)
	outputDeploy(t, deployList2)

	deployList3, err := handler.ListByNamespace(handler.namespace)
	myerr(t, "ListByNamespace", err)
	outputDeploy(t, deployList3)

	deployList4, err := handler.ListAll()
	myerr(t, "ListAll", err)
	outputDeploy(t, deployList4)

	handler.Delete(filename)
	handler.Delete(filename2)
}

func testWatchDeployment(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	deploy, err := handler.Apply(filename)
	if err != nil {
		t.Fatal(err)
	}

	addFunc := func(x interface{}) { t.Log("added deployment.") }
	modifyFunc := func(x interface{}) { t.Log("modified deployment.") }
	deleteFunc := func(x interface{}) { t.Log("deleted deployment.") }

	{
		t.Log("=== Watch")
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		timer := time.NewTimer(time.Second * 10)

		go func(ctx context.Context) {
			err = handler.Watch(deploy.Name, addFunc, modifyFunc, deleteFunc, nil)
			myerr(t, "Watch", err)
		}(ctx)
		go func(ctx context.Context) {
			time.Sleep(time.Second * 2)
			handler.Delete(name1)
			time.Sleep(time.Second)
			handler.Apply(filename)
		}(ctx)

		<-timer.C
		cancel()
	}
	{
		t.Log("=== WatchByName")
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		timer := time.NewTimer(time.Second * 10)

		go func(ctx context.Context) {
			err = handler.WatchByName(deploy.Name, addFunc, modifyFunc, deleteFunc, nil)
			myerr(t, "Watch", err)
		}(ctx)
		go func(ctx context.Context) {
			time.Sleep(time.Second * 2)
			handler.Delete(name1)
			time.Sleep(time.Second)
			handler.Apply(filename)
		}(ctx)

		<-timer.C
		cancel()
	}
	{
		t.Log("=== WatchByLabele")
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		timer := time.NewTimer(time.Second * 10)

		go func(ctx context.Context) {
			err = handler.WatchByLabel(label, addFunc, modifyFunc, deleteFunc, nil)
			myerr(t, "Watch", err)
		}(ctx)
		go func(ctx context.Context) {
			time.Sleep(time.Second * 2)
			handler.Delete(name1)
			time.Sleep(time.Second)
			handler.Apply(filename)
		}(ctx)

		<-timer.C
		cancel()
	}

}

func testDeploymentTools(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	deploy, err := handler.Apply(filename)
	if err != nil {
		t.Fatal(err)
	}
	// test IsReady, WaitReady
	t.Logf("deployment/%s is ready: %t", deploy.Name, handler.IsReady(deploy.Name))
	handler.WaitReady(deploy.Name)
	t.Logf("deployment/%s is ready: %t", deploy.Name, handler.IsReady(deploy.Name))

	// test GetRS
	rsList, err := handler.GetRS(name1)
	myerr(t, "GetRS", err)
	outputRS(t, rsList)

	// test GetPods
	podList, err := handler.GetPods(name1)
	myerr(t, "GetPods", err)
	outputPods(t, podList)
	handler.Delete(name1)

	// test GetPVC
	name := "nginx-deploy"
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
func outputDeploy(t *testing.T, deployList *appsv1.DeploymentList) {
	var dl []string
	for _, deploy := range deployList.Items {
		dl = append(dl, deploy.Name)
	}
	t.Log(dl)
}
func outputRS(t *testing.T, rsList []appsv1.ReplicaSet) {
	var rl []string
	for _, r := range rsList {
		rl = append(rl, r.Name)
	}
	t.Log(rl)
}
func outputPods(t *testing.T, podList []corev1.Pod) {
	var pl []string
	for _, p := range podList {
		pl = append(pl, p.Name)
	}
	t.Log(pl)
}
