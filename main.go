/*
GoZer is a prototypical database backend that was used to evaluate new ways of interaction with open access content on the go.
*/
package main

import (
	"os"
	"os/signal"
	"syscall"
)

import (
	"github.com/fzi-forschungszentrum-informatik/gozer/config"
	"github.com/fzi-forschungszentrum-informatik/gozer/storage"
	"github.com/fzi-forschungszentrum-informatik/gozer/storage/ledger"
	"github.com/fzi-forschungszentrum-informatik/gozer/webapi"
)

// waitForTerminateSignal waits until the GoZer service is stopped, either by an interrupt or terminate signal.
func waitForTerminateSignal() {

	// Set up a buffered channel for shutdown requests.
	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Wait until termination or interrupt signal is raised.
	<-sigs
}

// main runs the GoZer service until an interrupt or terminate signal is raised.
func main() {

	var webapi webapi.Service

	conf := config.LoadFromFile()
	storage := storage.Open(&conf.Storage)
	ledger := ledger.Open(&conf.Ledger)

	go webapi.Run(&conf.WebAPI, storage, ledger)

	waitForTerminateSignal()

	webapi.Shutdown()
	if ledger != nil {
		ledger.Close()
	}
	storage.Close()
}
