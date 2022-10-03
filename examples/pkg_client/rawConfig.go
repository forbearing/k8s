package main

import (
	"log"

	"github.com/forbearing/k8s"
)

func RawConfig() {
	rawConfig, err := k8s.RawConfig("")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("current context: ", rawConfig.CurrentContext)
}
