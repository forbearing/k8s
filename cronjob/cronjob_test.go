package cronjob

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var (
	ctx, cancel = context.WithTimeout(context.Background(), time.Second*10)
	namespace   = "test"
	kubeconfig  = filepath.Join(os.Getenv("HOME"), ".kube/config")
	path        = "../testData/examples/cronjob.yaml"
)

func TestCronjob(t *testing.T) {
	defer cancel()

	t.Run("Create CronJob", testCreateCronjob)
	t.Run("Update CronJob", testUpdateCronjob)
	t.Run("Apply Cronjob", testApplyCronjob)
	t.Run("Delete Cronjob", testDeleteCronjob)
	t.Run("Get Cronjob", testGetCronjob)
	t.Run("List Cronjob", testListCronjob)
	t.Run("Watch Cronjob", testWatchCronjob)
	t.Run("Cronjob Tools", testCronjobTools)
}

func testCreateCronjob(t *testing.T) {
	handler, err := New(ctx, namespace, kubeconfig)
	if err != nil {
		t.Fatal(err)
	}
	handler.DeleteFromFile(path)

	_, err = handler.Create(path)
	myerror(t, "Create", err)
	handler.DeleteFromFile(path)

	_, err = handler.CreateFromFile(path)
	myerror(t, "CreateFromFile", err)
	handler.DeleteFromFile(path)

	var data []byte
	handler.DeleteFromFile(path)
	if data, err = ioutil.ReadFile(path); err != nil {
		t.Error("ioutil.ReadFile error:", err)
	}
	_, err = handler.CreateFromBytes(data)
	myerror(t, "CreateFromBytes", err)
	handler.DeleteFromFile(path)

	name := "mycj-raw"
	rawData := map[string]interface{}{
		"apiVersion": "patch/v1",
		"kind":       "CronJob",
		"metadata": map[string]interface{}{
			"name": name,
			"labels": map[string]interface{}{
				"type": "cronjob",
			},
		},
		"spec": map[string]interface{}{
			"schedule": "*/1 * * * *",
			"jobTemplate": map[string]interface{}{
				"spec": map[string]interface{}{
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"containers": []map[string]interface{}{
								{
									"name":  "hello",
									"image": "busybox",
									"args:": []string{"/bin/sh", "-c", "date; echo hello kubernetes."},
								},
							},
							"restartPolicy": "OnFailure",
						},
					},
				},
			},
		},
	}
	_, err = handler.CreateFromRaw(rawData)
	myerror(t, "CreateFromRaw", err)
	handler.Delete(name)
}

func testUpdateCronjob(t *testing.T) {}
func testApplyCronjob(t *testing.T)  {}
func testDeleteCronjob(t *testing.T) {}
func testGetCronjob(t *testing.T)    {}
func testListCronjob(t *testing.T)   {}
func testWatchCronjob(t *testing.T)  {}
func testCronjobTools(t *testing.T)  {}

func myerror(t *testing.T, name string, err error) {
	if err != nil {
		t.Errorf("%s: %v", name, err)
	} else {
		t.Logf("%s success.", name)
	}
}
