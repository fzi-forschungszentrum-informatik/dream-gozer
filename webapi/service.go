package webapi

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"
)

import (
	"github.com/fzi-forschungszentrum-informatik/gozer/config"
	"github.com/fzi-forschungszentrum-informatik/gozer/storage"
	"github.com/fzi-forschungszentrum-informatik/gozer/storage/ledger"
)

// A Service object represents a Web service and its running state.
// The running state allows to wait for graceful shutdown.
type Service struct {
	server http.Server
	down   chan bool // signal for successful shutdown
}

// Starts the Web service, using the parameters specified in the global configuration.
// All request handler functions operate on the specified storage and ledger.
func (ws *Service) Run(conf *config.WebAPIConfiguration, st *storage.Storage, ledger *ledger.Ledger) {

	ws.down = make(chan bool, 1)

	ws.server.Addr = conf.Interface + ":" + strconv.Itoa(conf.Port)
	ws.server.Handler = newRouter(conf, st, ledger)
	ws.server.ReadTimeout = 15 * time.Second
	ws.server.WriteTimeout = 15 * time.Second

	log.Printf("Ploc web service is now listening on '%s'.", ws.server.Addr)

	if err := ws.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Starting ploc web service has failed. %s", err)
	}

	// Raise signal that service has shut down.
	ws.down <- true
}

// Initiates a gracefull shutdown of the Web service.
func (ws *Service) Shutdown() {

	err := ws.server.Shutdown(context.Background())

	<-ws.down

	if err != nil {
		log.Printf("Stopping ploc web service has failed. %v", err)
	} else {
		log.Print("Stopping ploc web service was successful.")
	}
}
