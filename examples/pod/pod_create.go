package main

import (
	"io/ioutil"

	"github.com/forbearing/k8s/pod"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Pod_Create() {
	handler, err := pod.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)

	// 1. create pod from filename
	pod, err := handler.Create(filename)
	checkErr("create pod from filename", pod.Name, err)
	handler.Delete(name)

	// 2. create pod from bytes
	var data []byte
	if data, err = ioutil.ReadFile(filename); err != nil {
		panic(err)
	}
	wait(handler, name)
	pod2, err := handler.Create(data)
	checkErr("create pod from bytes", pod2.Name, err)
	handler.Delete(name)

	// 3. create pod from *corev1.Pod
	wait(handler, name)
	pod3, err := handler.Create(pod2)
	checkErr("create pod from *corev1.pod", pod3.Name, err)
	handler.Delete(name)

	// 4. create pod from corev1.Pod
	wait(handler, name)
	pod4, err := handler.Create(*pod3)
	checkErr("create pod from corev1.Pod", pod4.Name, err)
	handler.Delete(name)

	// 5. create pod from runtime.Object
	wait(handler, name)
	pod5, err := handler.Create(runtime.Object(pod4))
	checkErr("create pod from runtime.Object", pod5.Name, err)
	handler.Delete(name)

	// 6. create pod from *unstructured.Unstructured
	pod6, err := handler.Create(&unstructured.Unstructured{Object: LogPodData})
	checkErr("create pod from *unstructured.Unstructured", pod6.Name, err)
	handler.Delete(LogPodName)

	// 7. create pod from unstructured.Unstructured
	wait(handler, LogPodName)
	pod7, err := handler.Create(unstructured.Unstructured{Object: LogPodData})
	checkErr("create pod from unstructured.Unstructured", pod7.Name, err)
	handler.Delete(LogPodName)

	// 8. create pod from map[string]interface{}
	wait(handler, LogPodName)
	pod8, err := handler.Create(LogPodData)
	checkErr("create pod from map[string]interface{}", pod8.Name, err)
	handler.Delete(LogPodName)

	// Output:

	//2022/08/08 16:41:28 create pod from filename success: mypod.
	//2022/08/08 16:41:28 create pod from bytes success: mypod.
	//2022/08/08 16:41:38 create pod from *corev1.pod success: mypod.
	//2022/08/08 16:41:38 create pod from corev1.Pod success: mypod.
	//2022/08/08 16:41:38 create pod from runtime.Object success: mypod.
	//2022/08/08 16:41:38 create pod from *unstructured.Unstructured success: nginx-logs.
	//2022/08/08 16:41:45 create pod from unstructured.Unstructured success: nginx-logs.
	//2022/08/08 16:41:45 create pod from map[string]interface{} success: nginx-logs.
}
