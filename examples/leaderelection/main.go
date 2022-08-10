package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/forbearing/k8s"
	"github.com/forbearing/k8s/utils/leaderelection"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	leaderElectionID        = "my-test-id"
	leaderElectionName      = "my-test-name"
	leaderElectionNamespace = "test"
)

func main() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Info("Received termination, signaling shutdown")
		cancel()
	}()

	onStartedLeading := func() {
		log.Info("Starting myserver")
		select {}
	}
	onStoppedLeading := func() {
		// start cleanup
		log.Infof("leader lost: %s", leaderElectionID)
		os.Exit(0)
	}
	onNewLeader := func(identity string) {
		// we're notified when new leader elected
		if identity == leaderElectionID {
			// I just got the lock
			return
		}
		log.Infof("new leader elected: %s", identity)
	}

	options := leaderelection.Options{
		LeaderElectionID:        leaderElectionID,
		LeaderElectionName:      leaderElectionName,
		LeaderElectionNamespace: leaderElectionNamespace,
		LeaseDuration:           time.Second * 15,
		RenewDeadline:           time.Second * 10,
		RetryPeriod:             time.Second * 2,
		OnStartedLeading:        onStartedLeading,
		OnStoppedLeading:        onStoppedLeading,
		OnNewLeader:             onNewLeader,
	}

	leaderelection.RunOrDie(ctx, k8s.RESTConfigOrDie(clientcmd.RecommendedHomeFile), options)
}
