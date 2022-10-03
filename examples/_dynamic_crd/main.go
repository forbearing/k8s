package main

import (
	"context"
	"log"

	storagev1alpha1 "github.com/forbearing/horus-operator/apis/storage/v1alpha1"
	"github.com/forbearing/k8s/dynamic"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	handler := dynamic.NewOrDie(context.TODO(), "", "default")
	gvk := schema.GroupVersionKind{
		Group:   storagev1alpha1.GroupVersion.Group,
		Version: storagev1alpha1.GroupVersion.Version,
		Kind:    "Backup",
	}
	_ = gvk

	var (
		err         error
		name        = "elassandra"
		filename    = "./iot-elassandra.yaml"
		unstructObj *unstructured.Unstructured
		backupObj   *storagev1alpha1.Backup
	)

	// 1.create from filename
	handler.DeleteFromFile(filename)
	if unstructObj, err = handler.Create(filename); err != nil {
		log.Fatalf("create failed: %v", err)
	}
	log.Printf("Successfully create %s from filename\n", convertUnstructObj(unstructObj).GetName())
	if err = handler.DeleteFromFile(filename); err != nil {
		log.Fatalf("delete failed: %v", err)
	}
	log.Printf("Successfully delete %s from filename", name)
	// 1.create from backup object
	backupObj = convertUnstructObj(unstructObj)
	if unstructObj, err = handler.Create(backupObj); err != nil {
		log.Fatalf("create failed: %v", err)
	}
	log.Printf("Successfully create %s from backup object\n", convertUnstructObj(unstructObj).GetName())
	if err = handler.WithGVK(gvk).Delete(name); err != nil {
		log.Fatalf("delete failed: %v", err)
	}
	log.Printf("Successfully delete %s from backup object name", name)

	// 2.update from backup object
	handler.Create(filename)
	uid := unstructObj.GetUID()
	resourceVersion := unstructObj.GetResourceVersion()
	unstructObj.SetResourceVersion(resourceVersion)
	unstructObj.SetUID(uid)
	//logrus.Info(unstructObj.GetResourceVersion())
	if unstructObj, err = handler.Update(unstructObj); err != nil {
		log.Fatalf("update failed: %v", err)
	}
	log.Printf("Successfully update %s\n", convertUnstructObj(unstructObj).GetName())
	handler.WithGVK(gvk).Delete(name)

	// 3.apply from filename
	handler.WithGVK(gvk).Delete(name)
	if unstructObj, err = handler.Apply(filename); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully apply %s from filename(not exist)\n", convertUnstructObj(unstructObj).GetName())
	// test failed
	//if unstructObj, err = handler.Apply(filename); err != nil {
	//    log.Fatal(err)
	//}
	//log.Printf("Successfully apply %s from filename(already exist)\n", convertUnstructObj(unstructObj).GetName())

	//unstructObj, err = handler.WithGVK(gvk).Get("iot-elassandra")
	//if err != nil {
	//    log.Println(unstructObj.GetName())
	//}
	handler.WithGVK(gvk).Delete(name)
	backupObj = convertUnstructObj(unstructObj)
	if unstructObj, err = handler.Apply(backupObj); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully apply %s from backup object(not exist)\n", convertUnstructObj(unstructObj).GetName())
	backupObj = convertUnstructObj(unstructObj)
	logrus.Info(backupObj)
	if unstructObj, err = handler.Apply(backupObj); err != nil {
		log.Fatal(err)
	}
	log.Printf("Successfully apply %s from backup object(already exist)\n", convertUnstructObj(unstructObj).GetName())
}

func convertUnstructObj(unstructObj *unstructured.Unstructured) *storagev1alpha1.Backup {
	backupObj := &storagev1alpha1.Backup{}
	if err := runtime.DefaultUnstructuredConverter.FromUnstructured(unstructObj.UnstructuredContent(), backupObj); err != nil {
		log.Fatal(err)
	}
	return backupObj
}
