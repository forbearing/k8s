package leaderelection

import (
	"time"
)

// This package draws heavily from the controller-runtime's leaderelection package
// (https://github.com/kubernetes-sigs/controller-runtime/tree/v0.12.3/pkg/leaderelection)
// but has some changes to bring it in line with my package style

const inClusterNamespacePath = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

// Options provides the requested configuration to create a new resource lock.
type Options struct {
	// LeaderElectionID is the unique string identifying a lease holder across
	// all participants in an election.
	LeaderElectionID string

	// LeaderElectionNamespace determines the namespace in which the leader
	// election resource will be created.
	LeaderElectionNamespace string

	// LeaderElectionName determines the name of the resource that leader election
	// will use for holding the leader lock.
	LeaderElectionName string

	// LeaseDuration is the duration that non-leader candidates will
	// wait to force acquire leadership. This is measured against time of
	// last observed ack.
	//
	// A client needs to wait a full LeaseDuration without observing a change to
	// the record before it can attempt to take over. When all clients are
	// shutdown and a new set of clients are started with different names against
	// the same leader record, they must wait the full LeaseDuration before
	// attempting to acquire the lease. Thus LeaseDuration should be as short as
	// possible (within your tolerance for clock skew rate) to avoid a possible
	// long waits in the scenario.
	//
	// Core clients default this value to 15 seconds.
	LeaseDuration time.Duration

	// RenewDeadline is the duration that the acting master will retry
	// refreshing leadership before giving up.
	//
	// Core clients default this value to 10 seconds.
	RenewDeadline time.Duration

	// RetryPeriod is the duration the LeaderElector clients should wait
	// between tries of actions.
	//
	// Core clients default this value to 2 seconds.
	RetryPeriod time.Duration

	// OnStartedLeading is called when a LeaderElector client starts leading
	OnStartedLeading func()

	// OnStoppedLeading is called when a LeaderElector client stops leading
	OnStoppedLeading func()

	// OnNewLeader is called when the client observes a leader that is
	// not the previously observed leader. This includes the first observed
	// leader when the client starts.
	OnNewLeader func(identity string)
}
