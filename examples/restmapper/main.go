package main

import (
	"log"
)

func main() {
	//FindGVK()
	FindGVR()
	//Is_Namespaced()
}

func checkErr(name string, val interface{}, err error) {
	if err != nil {
		log.Printf("%s failed: %v\n", name, err)
	} else {
		log.Printf("%s success: %v\n", name, val)
	}
}
