package config

import (
	"flag"
	"fmt"
	"log"
)

import (
	"github.com/BurntSushi/toml"
)

const (
	defaultConfigFilename = "gozer.conf"
)

// Defines the global configuration parameters of GoZer's Web API.
// The interface defines the network device GoZer is using for communication, e.g. "192.168.222.1" (specific interface)
// or "0.0.0.0" (all interfaces). The port specifies where GoZer is listening for HTTP requests. In addition a local path
// to ploc APK file can be specified, which allows to download the APK from GoZer directly.
type WebAPIConfiguration struct {
	Interface string `toml:"interface"`
	Port      int    `toml:"port"`
	PlocAPK   string `toml:"ploc_apk"`
}

// Defines the global configuration parameters for GoZer's SQLite database, which is the local path to the SQLite database file.
type StorageConfiguration struct {
	DBFilename string `toml:"db_filename"`
}

// Defines the global configuration parameters for GoZer's Ethereum connector.
// The use of Ethereum is optional an can be disabled.
// The RPC client specifies the URL to a (public) Ethereum node.
// The contract address defines under which address the smartcontract for open feedback was deployed.
// The private key specifies the identity / wallet that GoZer uses to publish open feedback.
type LedgerConfiguration struct {
	Enable          bool   `toml:"enable"`
	RPCClient       string `toml:"rpc_client"`
	ContractAddress string `toml:"contract_address"`
	PrivateKey      string `toml:"private_key"`
}

// Defines the global configuration of the GoZer service.
// The configuration consists of parameters for the Web API, the database and the Ethereum ledger.
type Configuration struct {
	WebAPI  WebAPIConfiguration  `toml:"webapi"`
	Storage StorageConfiguration `toml:"storage"`
	Ledger  LedgerConfiguration  `toml:"ledger"`
}

// DefaultConfiguration returns a default configuration, that can be used e.g. for testing.
// The ledger is disabled per default.
func DefaultConfiguration() *Configuration {

	var conf Configuration

	conf.WebAPI.Interface = "0.0.0.0"
	conf.WebAPI.Port = 8080
	conf.WebAPI.PlocAPK = "ploc.apk"

	conf.Storage.DBFilename = "storage.db"

	conf.Ledger.Enable = false
	conf.Ledger.RPCClient = "http://127.0.0.1:7545"                                             // e.g. for testing with Ganache
	conf.Ledger.ContractAddress = "0xc8B381DCCAE278F809DB5e0b7B2EfA8c716270d7"                  // dummy, not a real one
	conf.Ledger.PrivateKey = "0fc142ddbe063614c3cab903fbc1516a5ab663d1fa8bcfb46f867df4bd5c03fe" // dummy, not a real one

	return &conf
}

// LoadFromFile loads a configuration from a TOML-based file. The path to the file can either be specified as command
// line parameter via the '-f' option or is assumed to be 'gozer.conf'.
func LoadFromFile() *Configuration {

	var f string

	// Create a new configuration with some default values.
	conf := DefaultConfiguration()

	// Read config filename specified as command line parameter ('-f name') or the default name.
	flag.StringVar(&f, "f", defaultConfigFilename, fmt.Sprintf("Specifies a configuration file. Default is '%s'.", defaultConfigFilename))
	flag.Parse()

	// Use the BurntSushi library to load and parse TOML configuration file.
	if _, err := toml.DecodeFile(f, conf); err != nil {
		log.Fatalf("Reading configuration '%s' has failed. %s", f, err)
	}

	log.Printf("Reading configuration '%s' was successfull.", f)

	return conf
}
