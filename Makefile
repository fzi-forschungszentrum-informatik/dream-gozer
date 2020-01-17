GOZER_URI := github.com/fzi-forschungszentrum-informatik/gozer

GOFILES = $(shell find . -name "*.go" -type f) 

gozer: $(GOFILES)
	go build --tags "fts5" -i -o gozer .

## Generate Go code from SQL schemas and Solidity contract

./storage/schema.go: ./storage/schema.sql
	go generate $(GOZER_URI)/storage/

./storage/testdb.go: ./storage/testdb.sql
	go generate $(GOZER_URI)/storage/

./storage/ledger/open_feedback.go: ./storage/ledger/open_feedback.sol
	abigen --sol=./storage/ledger/open_feedback.sol --pkg=ledger --out=./storage/ledger/open_feedback.go

## Testing

test:
	go test -cover ./... --tags "fts5"
	# go test github.com/fzi-forschungszentrum-informatik/gozer/webapi -run TestFeedbackFeed --tags "fts5"

## Documentation

doc:
	# go to 'http://localhost:8080/pkg/github.com/fzi-forschungszentrum-informatik/gozer/'
	godoc -http=:8080 -notes="BUG|TODO"

## Cleanup

clean:
	rm -f gozer 

## Dependencies

install_packages:
	go get -u "github.com/google/uuid" "github.com/mattn/go-sqlite3" "github.com/gorilla/mux" "golang.org/x/crypto/bcrypt" "github.com/BurntSushi/toml" "github.com/ethereum/go-ethereum"
	# go get -u "github.com/ethereum/go-ethereum/..."
	# 2019-11-26: Broken package-dependency in go-ethereum (influxdb client missing), but in a part we don't need (abigen) for GoZer.
	go get -u -u "github.com/ethereum/go-ethereum" "github.com/davecgh/go-spew/spew" "github.com/deckarep/golang-set" "github.com/edsrzf/mmap-go" "github.com/gballet/go-libpcsclite" "github.com/golang/protobuf/proto" "github.com/golang/protobuf/protoc-gen-go/descriptor" "github.com/gorilla/websocket" "github.com/hashicorp/golang-lru" "github.com/hashicorp/golang-lru/simplelru" "github.com/huin/goupnp" "github.com/huin/goupnp/dcps/internetgateway1" "github.com/huin/goupnp/dcps/internetgateway2" "github.com/jackpal/go-nat-pmp" "github.com/karalabe/usb" "github.com/olekukonko/tablewriter" "github.com/pborman/uuid" "github.com/prometheus/tsdb/fileutil" "github.com/rjeczalik/notify" "github.com/rs/cors" "github.com/status-im/keycard-go/derivationpath" "github.com/syndtr/goleveldb/leveldb" "github.com/syndtr/goleveldb/leveldb/errors" "github.com/syndtr/goleveldb/leveldb/filter" "github.com/syndtr/goleveldb/leveldb/iterator" "github.com/syndtr/goleveldb/leveldb/opt" "github.com/syndtr/goleveldb/leveldb/storage" "github.com/syndtr/goleveldb/leveldb/util" "github.com/tyler-smith/go-bip39" "github.com/wsddn/go-ecdh"

.PHONY: test clean install_packages
