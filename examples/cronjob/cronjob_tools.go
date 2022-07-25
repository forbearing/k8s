package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/cronjob"
	batchv1 "k8s.io/api/batch/v1"
)

func Cronjob_Tools() {
	handler, err := cronjob.New(ctx, kubeconfig, namespace)
	if err != nil {
		panic(err)
	}
	//defer cleanup(handler)

	handler.Apply(filename)
	cj, err := handler.Get(name)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get Cronjob By Name")
		begin := time.Now()

		jobs, err := handler.GetJobs(name)
		checkErr("GetJobs", printJobs(jobs), err)
		numActive, err := handler.NumActive(name)
		checkErr("NumActive", numActive, err)
		lastScheduled, err := handler.DurationOfLastScheduled(name)
		checkErr("DurationOfLastScheduled", lastScheduled, err)
		completed, err := handler.DurationOfCompleted(name)
		checkErr("DurationOfCompleted", completed, err)
		schedule, err := handler.GetSchedule(name)
		checkErr("GetSchedule", schedule, err)
		suspend, err := handler.IsSuspend(name)
		checkErr("IsSuspend", suspend, err)
		age, err := handler.GetAge(name)
		checkErr("GetAge", age, err)
		containers, err := handler.GetContainers(name)
		checkErr("GetContainers", containers, err)
		images, err := handler.GetImages(name)
		checkErr("GetImages", images, err)

		end := time.Now()
		log.Println("===== Get Cronjob By Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get Cronjob By Object")
		begin := time.Now()

		jobs, err := handler.GetJobs(cj)
		checkErr("GetJobs", printJobs(jobs), err)
		numActive, err := handler.NumActive(cj)
		checkErr("NumActive", numActive, err)
		lastScheduled, err := handler.DurationOfLastScheduled(cj)
		checkErr("DurationOfLastScheduled", lastScheduled, err)
		completed, err := handler.DurationOfCompleted(cj)
		checkErr("DurationOfCompleted", completed, err)
		schedule, err := handler.GetSchedule(cj)
		checkErr("GetSchedule", schedule, err)
		suspend, err := handler.IsSuspend(cj)
		checkErr("IsSuspend", suspend, err)
		age, err := handler.GetAge(cj)
		checkErr("GetAge", age, err)
		containers, err := handler.GetContainers(cj)
		checkErr("GetContainers", containers, err)
		images, err := handler.GetImages(cj)
		checkErr("GetImages", images, err)

		end := time.Now()
		log.Println("===== Get Cronjob By Object Cost Time:", end.Sub(begin))
	}

	getByName()
	getByObj()

	// Output:

	//2022/07/11 14:09:30 ===== Get Cronjob By Name
	//2022/07/11 14:09:30 GetJobs success: [mycj-27625326 mycj-27625327 mycj-27625328 mycj-27625329].
	//2022/07/11 14:09:30 NumActive success: 1.
	//2022/07/11 14:09:30 DurationOfLastScheduled success: 30.44301s.
	//2022/07/11 14:09:30 DurationOfCompleted failed: the last time the job successfully completed not found
	//2022/07/11 14:09:30 GetSchedule success: */1 * * * *.
	//2022/07/11 14:09:30 IsSuspend success: false.
	//2022/07/11 14:09:30 GetAge success: 71h34m28.558962s.
	//2022/07/11 14:09:30 GetContainers success: [hello].
	//2022/07/11 14:09:30 GetImages success: [busybox].
	//2022/07/11 14:09:30 ===== Get Cronjob By Name Cost Time: 559.135424ms
	//2022/07/11 14:09:30
	//2022/07/11 14:09:30 ===== Get Cronjob By Object
	//2022/07/11 14:09:31 GetJobs success: [mycj-27625326 mycj-27625327 mycj-27625328 mycj-27625329].
	//2022/07/11 14:09:31 NumActive success: 1.
	//2022/07/11 14:09:31 DurationOfLastScheduled success: 31.128522s.
	//2022/07/11 14:09:31 DurationOfCompleted failed: the last time the job successfully completed not found
	//2022/07/11 14:09:31 GetSchedule success: */1 * * * *.
	//2022/07/11 14:09:31 IsSuspend success: false.
	//2022/07/11 14:09:31 GetAge success: 71h34m29.128578s.
	//2022/07/11 14:09:31 GetContainers success: [hello].
	//2022/07/11 14:09:31 GetImages success: [busybox].
	//2022/07/11 14:09:31 ===== Get Cronjob By Object Cost Time: 191.884047ms
}

func printJobs(jobList []batchv1.Job) []string {
	var cl []string
	for _, job := range jobList {
		cl = append(cl, job.Name)
	}
	return cl
}
