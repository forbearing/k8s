package main

import (
	"log"
	"time"

	"github.com/forbearing/k8s/job"
)

func Job_Tools() {
	handler, err := job.New(ctx, namespace, kubeconfig)
	if err != nil {
		panic(err)
	}
	handler.Apply(filename)
	handler.Apply(filename2)
	j1, err := handler.Get(name)
	if err != nil {
		panic(err)
	}
	j2, err := handler.Get(name2)
	if err != nil {
		panic(err)
	}

	getByName := func() {
		log.Println("===== Get From Job Name")
		begin := time.Now()

		startTime1, err := handler.DurationOfStarted(name)
		checkErr("DurationOfStarted", err)
		startTime2, err := handler.DurationOfStarted(name2)
		checkErr("DurationOfStarted", err)
		log.Println(startTime1)
		log.Println(startTime2)

		compTime1, err := handler.DurationOfCompleted(name)
		checkErr("DurationOfCompleted", err)
		compTime2, err := handler.DurationOfCompleted(name2)
		checkErr("DurationOfCompleted", err)
		log.Println(compTime1)
		log.Println(compTime2)

		numActive1, err := handler.NumActive(name)
		checkErr("NumActive", err)
		numActive2, err := handler.NumActive(name2)
		checkErr("NumActive", err)
		log.Println(numActive1)
		log.Println(numActive2)

		numSucc1, err := handler.NumSucceeded(name)
		checkErr("NumSucceeded", err)
		numSucc2, err := handler.NumSucceeded(name2)
		checkErr("NumSucceeded", err)
		log.Println(numSucc1)
		log.Println(numSucc2)

		numFailed1, err := handler.NumFailed(name)
		checkErr("NumFailed", err)
		numFailed2, err := handler.NumFailed(name2)
		checkErr("NumFailed", err)
		log.Println(numFailed1)
		log.Println(numFailed2)

		age1, err := handler.GetAge(name)
		checkErr("GetAge", err)
		age2, err := handler.GetAge(name2)
		checkErr("GetAge", err)
		log.Println(age1)
		log.Println(age2)

		end := time.Now()
		log.Println("===== Get From Job Name Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByObj := func() {
		log.Println("===== Get From Job Object")
		begin := time.Now()

		startTime1, err := handler.DurationOfStarted(j1)
		checkErr("DurationOfStarted", err)
		startTime2, err := handler.DurationOfStarted(j2)
		checkErr("DurationOfStarted", err)
		log.Println(startTime1)
		log.Println(startTime2)

		compTime1, err := handler.DurationOfCompleted(j1)
		checkErr("DurationOfCompleted", err)
		compTime2, err := handler.DurationOfCompleted(j2)
		checkErr("DurationOfCompleted", err)
		log.Println(compTime1)
		log.Println(compTime2)

		numActive1, err := handler.NumActive(j1)
		checkErr("NumActive", err)
		numActive2, err := handler.NumActive(j2)
		checkErr("NumActive", err)
		log.Println(numActive1)
		log.Println(numActive2)

		numSucc1, err := handler.NumSucceeded(j1)
		checkErr("NumSucceeded", err)
		numSucc2, err := handler.NumSucceeded(j2)
		checkErr("NumSucceeded", err)
		log.Println(numSucc1)
		log.Println(numSucc2)

		numFailed1, err := handler.GetAge(j1)
		checkErr("NumFailed", err)
		numFailed2, err := handler.NumFailed(j2)
		checkErr("NumFailed", err)
		log.Println(numFailed1)
		log.Println(numFailed2)

		age1, err := handler.GetAge(j1)
		checkErr("GetAge", err)
		age2, err := handler.GetAge(j2)
		checkErr("GetAge", err)
		log.Println(age1)
		log.Println(age2)

		end := time.Now()
		log.Println("===== Get From Job Object Cost Time:", end.Sub(begin))
		log.Println()
	}

	getByName()
	getByObj()

	// Output

	//2022/07/08 13:52:22 ===== Get From Job Name
	//2022/07/08 13:52:22 DurationOfStarted success.
	//2022/07/08 13:52:22 DurationOfStarted success.
	//2022/07/08 13:52:22 3h0m47.019039s
	//2022/07/08 13:52:22 2h30m23.024747s
	//2022/07/08 13:52:22 DurationOfCompleted success.
	//2022/07/08 13:52:22 DurationOfCompleted failed: completion time not found
	//2022/07/08 13:52:22 3h0m13.034342s
	//2022/07/08 13:52:22 0s
	//2022/07/08 13:52:22 NumActive success.
	//2022/07/08 13:52:22 NumActive success.
	//2022/07/08 13:52:22 0
	//2022/07/08 13:52:22 0
	//2022/07/08 13:52:22 NumSucceeded success.
	//2022/07/08 13:52:22 NumSucceeded success.
	//2022/07/08 13:52:22 1
	//2022/07/08 13:52:22 0
	//2022/07/08 13:52:22 NumFailed success.
	//2022/07/08 13:52:23 NumFailed success.
	//2022/07/08 13:52:23 0
	//2022/07/08 13:52:23 4
	//2022/07/08 13:52:23 GetAge success.
	//2022/07/08 13:52:23 GetAge success.
	//2022/07/08 13:52:23 3h0m48.341862s
	//2022/07/08 13:52:23 2h30m25.55851s
	//2022/07/08 13:52:23 ===== Get From Job Name Cost Time: 1.547452384s
	//2022/07/08 13:52:23
	//2022/07/08 13:52:23 ===== Get From Job Object
	//2022/07/08 13:52:23 DurationOfStarted success.
	//2022/07/08 13:52:23 DurationOfStarted success.
	//2022/07/08 13:52:23 3h0m48.558571s
	//2022/07/08 13:52:23 2h30m24.558576s
	//2022/07/08 13:52:23 DurationOfCompleted success.
	//2022/07/08 13:52:23 DurationOfCompleted failed: completion time not found
	//2022/07/08 13:52:23 3h0m14.558589s
	//2022/07/08 13:52:23 0s
	//2022/07/08 13:52:23 NumActive success.
	//2022/07/08 13:52:23 NumActive success.
	//2022/07/08 13:52:23 0
	//2022/07/08 13:52:23 0
	//2022/07/08 13:52:23 NumSucceeded success.
	//2022/07/08 13:52:23 NumSucceeded success.
	//2022/07/08 13:52:23 1
	//2022/07/08 13:52:23 0
	//2022/07/08 13:52:23 NumFailed success.
	//2022/07/08 13:52:23 NumFailed success.
	//2022/07/08 13:52:23 3h0m48.558693s
	//2022/07/08 13:52:23 4
	//2022/07/08 13:52:23 GetAge success.
	//2022/07/08 13:52:23 GetAge success.
	//2022/07/08 13:52:23 3h0m48.558712s
	//2022/07/08 13:52:23 2h30m25.558716s
	//2022/07/08 13:52:23 ===== Get From Job Object Cost Time: 158.902µs
	//2022/07/08 13:52:23
}
