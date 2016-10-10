package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/andreas-schroeder/kafka-health-check/check"
)

func main() {
	healthCheck := check.New(checkConfiguration)
	healthCheck.ParseCommandLineArguments()

	stop, awaitCheck := addShutdownHook()
	brokerUpdates, clusterUpdates, zookeeperUpdates := make(chan string, 3), make(chan string, 3), make(chan string, 3)
	go healthCheck.ServeHealth(brokerUpdates, clusterUpdates, zookeeperUpdates, stop)
	healthCheck.CheckHealth(brokerUpdates, clusterUpdates, stop)
	awaitCheck.Done()
}

func addShutdownHook() (chan struct{}, sync.WaitGroup) {
	stop := make(chan struct{})
	awaitCheck := sync.WaitGroup{}
	awaitCheck.Add(1)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)
	signal.Notify(shutdown, syscall.SIGTERM)
	go func() {
		for _ = range shutdown {
			close(stop)
			awaitCheck.Wait()
		}
	}()

	return stop, awaitCheck
}

var checkConfiguration = check.HealthCheckConfig{
	CheckTimeout:     100 * time.Millisecond,
	DataWaitInterval: 20 * time.Millisecond,
	MessageLength:    20,
}
