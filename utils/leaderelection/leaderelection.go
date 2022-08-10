package leaderelection

// This package draws heavily from the controller-runtime's leaderelection package
// (https://github.com/kubernetes-sigs/controller-runtime/tree/v0.12.3/pkg/leaderelection)
// but has some changes to bring it in line with my package style

import (
	"context"
	"errors"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

const inClusterNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

// RunOrDie
func RunOrDie(ctx context.Context, config *rest.Config, options Options) {
	lec, err := NewLeaderElectionConfig(config, options)
	if err != nil {
		panic(err)
	}
	leaderelection.RunOrDie(ctx, *lec)
}

// NewLeaderElectionConfig
func NewLeaderElectionConfig(config *rest.Config, options Options) (*leaderelection.LeaderElectionConfig, error) {
	lock, err := NewResourceLock(config, nil, options)
	if err != nil {
		return nil, err
	}

	return &leaderelection.LeaderElectionConfig{
		Lock: lock,
		// IMPORTANT: you MUST ensure that any code you have that
		// is protected by the lease must terminate **before**
		// you call cancel. Otherwise, you could have a background
		// loop still running and another process could
		// get elected before your background loop finished, violating
		// the stated goal of the lease.
		ReleaseOnCancel: true,
		LeaseDuration:   options.LeaseDuration,
		RenewDeadline:   options.RenewDeadline,
		RetryPeriod:     options.RetryPeriod,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) { options.OnStartedLeading() },
			OnStoppedLeading: options.OnStoppedLeading,
			OnNewLeader:      options.OnNewLeader,
		},
	}, nil
}

// NewResourceLock creates a resourcelock.Interface object with provided parameters.
func NewResourceLock(config *rest.Config, eventRecorder record.EventRecorder, options Options) (resourcelock.Interface, error) {
	// LeaderElectionName must be provided to prevent clashes
	if len(options.LeaderElectionName) == 0 {
		return nil, errors.New("LeaderElectionName must be configured")
	}

	// Set the default namespace(if running in cluster)
	if len(options.LeaderElectionNamespace) == 0 {
		var err error
		options.LeaderElectionNamespace, err = getInClusterNamespace()
		if err != nil {
			return nil, fmt.Errorf("unable to find leader election namespace: %w", err)
		}
	}

	// Set the leader id and must be unique.
	if len(options.LeaderElectionID) == 0 {
		var err error
		options.LeaderElectionID, err = os.Hostname()
		if err != nil {
			return nil, err
		}
		options.LeaderElectionID = options.LeaderElectionID + "_" + string(uuid.NewUUID())

	}

	// Construct clientset for leader election.
	rest.AddUserAgent(config, "leader-election")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return resourcelock.New(
		resourcelock.LeasesResourceLock,
		options.LeaderElectionNamespace,
		options.LeaderElectionName,
		clientset.CoreV1(),
		clientset.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity:      options.LeaderElectionID,
			EventRecorder: eventRecorder,
		})
}

// getInClusterNamespace will get the namespace of currently running election node(a k8s pod)
func getInClusterNamespace() (string, error) {
	// Check whether the namespace file exists.
	if _, err := os.Stat(inClusterNamespacePath); os.IsNotExist(err) {
		return "", errors.New("not running in-cluster, please specify LeaderElectionNamespace")
	} else if err != nil {
		return "", fmt.Errorf("error checking namespace file: %w", err)
	}

	// Load the namespace file and return its content.
	namespace, err := os.ReadFile(inClusterNamespacePath)
	if err != nil {
		return "", fmt.Errorf("error reading namespace file: %w", err)
	}
	return string(namespace), nil
}
