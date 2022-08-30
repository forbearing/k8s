package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

// SetupSignalContext registers for SIGTERM and SIGINT. A context is returned
// which is canceled on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
//
// It will panic when call SetupSignalContext twice.
func SetupSignalContext() context.Context {
	close(onlyOneSignalHandler) // panics when called twice

	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		cancel()
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}

// SetupSignalChannel registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
//
// It will panic when call SetupSignalChannel twice.
func SetupSignalChannel() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

// NewSignalContext registers for SIGTERM and SIGINT. A context is returned
// which is canceled on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
//
// You can call NewSignalContext() multiple time.
func NewSignalContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		cancel()
		//<-c
		//os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}

// NewSignalChannel registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
//
// You can call NewSignalChannel() multiple time.
func NewSignalChannel() (stopCh <-chan struct{}) {
	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		close(stop)
		//<-c
		//os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}
