package main

import (
	"time"

	"github.com/forbearing/k8s/pod"
)

func Pod() {
	handler, err := pod.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	handler.Apply(filename)
	handler.WaitReady(name)

	// execute success.
	command1 := []string{
		//"sh",
		//"-c",
		"hostname",
	}
	time.Sleep(time.Second * 5)
	command2 := []string{
		"sh",
		"-c",
		"hostname",
	}
	command3 := []string{
		"/bin/sh",
		"-c",
		"hostname",
	}
	command4 := []string{
		"/bin/bash",
		"-c",
		"hostname",
	}
	command5 := []string{
		"cat /etc/os-release",
	}
	command6 := []string{
		"sh",
		"-c",
		"cat /etc/os-release",
	}
	command7 := []string{
		"sh",
		"-c",
		"apt update; apt upgrade -y",
	}
	handler.Execute(name, "", command1) // execute success.
	handler.Execute(name, "", command2) // execute success.
	handler.Execute(name, "", command3) // execute success.
	handler.Execute(name, "", command4) // execute success.
	handler.Execute(name, "", command5) // execute failed.
	handler.Execute(name, "", command6) // execute success.
	handler.Execute(name, "", command7) // execute success, but may be cancelled by context timeout.
}
