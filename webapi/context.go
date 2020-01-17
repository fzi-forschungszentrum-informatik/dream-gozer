package webapi

import (
	"github.com/fzi-forschungszentrum-informatik/dream-gozer/config"
	"github.com/fzi-forschungszentrum-informatik/dream-gozer/storage"
	"github.com/fzi-forschungszentrum-informatik/dream-gozer/storage/ledger"
)

// Context defines a state of information, in which a HTTP request is interpreted.
// In GoZer this state is composed by the state of the database and the configuration file.
type Context struct {
	conf   *config.WebAPIConfiguration
	db     *storage.Storage
	ledger *ledger.Ledger
}

// newContext defines a new context object, consisting of global configuration information, a data storage and an Ethereum ledger.
func newContext(conf *config.WebAPIConfiguration, db *storage.Storage, ledger *ledger.Ledger) *Context {
	return &Context{conf: conf, db: db, ledger: ledger}
}
