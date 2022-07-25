package main

import (
	"fmt"

	"github.com/forbearing/k8s/pod"
)

func Pod_Logs() {
	handler, err := pod.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	//defer cleanup(handler)

	// create a pod
	_, err = handler.Apply(LogPodData)
	if err != nil {
		panic(err)
	}

	// wait pod ready.
	fmt.Printf("wait pod/%s ready\n", LogPodName)
	handler.WaitReady(LogPodName)

	// get pod logs from pod name.
	err = handler.Log(LogPodName, pod.DefaultLogOptions)
	checkErr("get pod logs from pod name", "", err)
	fmt.Println()

	// get pod logs from unstructured data.
	err = handler.Log(LogPodData, pod.DefaultLogOptions)
	checkErr("get pod logs from unstructured data", "", err)

	// output:

	//wait pod/nginx-logs ready
	///docker-entrypoint.sh: /docker-entrypoint.d/ is not empty, will attempt to perform configuration
	///docker-entrypoint.sh: Looking for shell scripts in /docker-entrypoint.d/
	///docker-entrypoint.sh: Launching /docker-entrypoint.d/10-listen-on-ipv6-by-default.sh
	//10-listen-on-ipv6-by-default.sh: info: Getting the checksum of /etc/nginx/conf.d/default.conf
	//10-listen-on-ipv6-by-default.sh: info: Enabled listen on IPv6 in /etc/nginx/conf.d/default.conf
	///docker-entrypoint.sh: Launching /docker-entrypoint.d/20-envsubst-on-templates.sh
	///docker-entrypoint.sh: Launching /docker-entrypoint.d/30-tune-worker-processes.sh
	///docker-entrypoint.sh: Configuration complete; ready for start up
	//2022/07/25 09:35:27 [notice] 1#1: using the "epoll" event method
	//2022/07/25 09:35:27 [notice] 1#1: nginx/1.21.5
	//2022/07/25 09:35:27 [notice] 1#1: built by gcc 10.2.1 20210110 (Debian 10.2.1-6)
	//2022/07/25 09:35:27 [notice] 1#1: OS: Linux 4.15.0-188-generic
	//2022/07/25 09:35:27 [notice] 1#1: getrlimit(RLIMIT_NOFILE): 1048576:1048576
	//2022/07/25 09:35:27 [notice] 1#1: start worker processes
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 32
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 33
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 34
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 35
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 36
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 37
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 38
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 39
	//2022/07/25 17:35:28 get pod logs from pod name success: .

	///docker-entrypoint.sh: /docker-entrypoint.d/ is not empty, will attempt to perform configuration
	///docker-entrypoint.sh: Looking for shell scripts in /docker-entrypoint.d/
	///docker-entrypoint.sh: Launching /docker-entrypoint.d/10-listen-on-ipv6-by-default.sh
	//10-listen-on-ipv6-by-default.sh: info: Getting the checksum of /etc/nginx/conf.d/default.conf
	//10-listen-on-ipv6-by-default.sh: info: Enabled listen on IPv6 in /etc/nginx/conf.d/default.conf
	///docker-entrypoint.sh: Launching /docker-entrypoint.d/20-envsubst-on-templates.sh
	///docker-entrypoint.sh: Launching /docker-entrypoint.d/30-tune-worker-processes.sh
	///docker-entrypoint.sh: Configuration complete; ready for start up
	//2022/07/25 09:35:27 [notice] 1#1: using the "epoll" event method
	//2022/07/25 09:35:27 [notice] 1#1: nginx/1.21.5
	//2022/07/25 09:35:27 [notice] 1#1: built by gcc 10.2.1 20210110 (Debian 10.2.1-6)
	//2022/07/25 09:35:27 [notice] 1#1: OS: Linux 4.15.0-188-generic
	//2022/07/25 09:35:27 [notice] 1#1: getrlimit(RLIMIT_NOFILE): 1048576:1048576
	//2022/07/25 09:35:27 [notice] 1#1: start worker processes
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 32
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 33
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 34
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 35
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 36
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 37
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 38
	//2022/07/25 09:35:27 [notice] 1#1: start worker process 39
	//2022/07/25 17:35:28 get pod logs from unstructured data success: .
}
