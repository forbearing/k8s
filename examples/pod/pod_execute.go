package main

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/forbearing/k8s/pod"
)

func Pod_Execute() {
	handler, err := pod.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	defer cleanup(handler)
	handler.Apply(filename)
	handler.WaitReady(name)

	command1 := []string{
		"hostname",
	}
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
	handler.Execute(name, "", command1)                                           // execute success.
	handler.ExecuteWithPty(name, "", command2, nil)                               // execute success.
	handler.ExecuteWithPty(name, "", command3, nil)                               // execute success.
	handler.ExecuteWithStream(name, "", command4, os.Stdin, os.Stdout, os.Stderr) // execute success.
	handler.Execute(name, "", command5)                                           // execute failed.
	handler.Execute(name, "", command6)                                           // execute success.
	handler.Execute(name, "", command7)                                           // execute success, but may be cancelled by context timeout.

	// Output:

	//mypod
	//mypod
	//mypod
	//mypod
	//OCI runtime exec failed: exec failed: unable to start container process: exec: "cat /etc/os-release": stat cat /etc/os-release: no such file or directory: unknown
	//PRETTY_NAME="Debian GNU/Linux 11 (bullseye)"
	//NAME="Debian GNU/Linux"
	//VERSION_ID="11"
	//VERSION="11 (bullseye)"
	//VERSION_CODENAME=bullseye
	//ID=debian
	//HOME_URL="https://www.debian.org/"
	//SUPPORT_URL="https://www.debian.org/support"
	//BUG_REPORT_URL="https://bugs.debian.org/"
	//Get:1 http://security.debian.org/debian-security bullseye-security InRelease [44.1 kB]
	//Get:2 http://deb.debian.org/debian bullseye InRelease [116 kB]
	//Get:3 http://security.debian.org/debian-security bullseye-security/main amd64 Packages [163 kB]
	//Get:4 http://deb.debian.org/debian bullseye-updates InRelease [39.4 kB]
	//Get:5 http://deb.debian.org/debian bullseye/main amd64 Packages [8182 kB]
	//26% [5 Packages 750 kB/8182 kB 9%]                           36.7 kB/s 3min 22s

	file, err := os.Create("/tmp/pod-exec")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buffer := &bytes.Buffer{}
	_ = file
	_ = buffer
	handler.ExecuteWithStream(name, "", command1, os.Stdin, os.Stdout, os.Stderr)
	handler.ExecuteWithStream(name, "", command1, os.Stdin, file, file)
	handler.ExecuteWithStream(name, "", command1, os.Stdin, buffer, buffer)
	fmt.Println(buffer.String())

	// Output
	//mypod
	//mypod

	podName := "backup-to-nfs-f678db6b6-jqsl2"
	resticPass := &bytes.Buffer{}
	resticPass.WriteString("mypass")
	if err := handler.WithNamespace("default").ExecuteWithStream(podName, "", []string{"restic", "--repo", "/tmp/restic-repo", "init"}, resticPass, os.Stdout, os.Stderr); err != nil {
		log.Fatal(err)
	}
}
