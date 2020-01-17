package ledger

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fzi-forschungszentrum-informatik/dream-gozer/config"
)

// Ledger defines an abstraction to the Ethereum blockchain, that includes the contract and the private
// key for the GoZer client node.
type Ledger struct {
	Client     *ethclient.Client
	PrivateKey *ecdsa.PrivateKey
	Contract   *OpenFeedback
}

// Open initializes the ledger accordingly to the provided global configuration and connects to the Ethereum blockchain.
func Open(conf *config.LedgerConfiguration) *Ledger {

	if !conf.Enable {
		return nil
	}

	client, err := ethclient.Dial(conf.RPCClient)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum network. %s", err)
	}

	contract, err := NewOpenFeedback(common.HexToAddress(conf.ContractAddress), client)
	if err != nil {
		log.Fatalf("Failed to initialize contract. %s", err)
	}

	// Eliptic Curve DSA
	privateKey, err := crypto.HexToECDSA(conf.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to convert private key. %s", err)
	}

	log.Printf("Connecting to Ethereum network via '%s' was successfull.", conf.RPCClient)

	return &Ledger{PrivateKey: privateKey, Contract: contract, Client: client}
}

// Close finishes open transaction and disconnects then from the Ethereum blockchain.
func (st *Ledger) Close() {

	if st.Client != nil {
		st.Client.Close()
		log.Print("Closing connection to Ethereum network was successfull.")
	}
}

// AddFeedback publishes a user's feedback about a publication to the Ethereum blockchain.
// In the blockchain a user is identified by its ORCiD and the publication by its bibliographic hash.
// The feedback consists of a binary flag about the quality of relevance, presentation and methodology.
// The idea is, that feedback about a scientific publication is public and accessible to anyone.
func (st *Ledger) AddFeedback(orcId string, bibHash string, relevance uint8, presentation uint8, methodology uint8) (err error) {

	binOrcId, err := orcIdToByteArray(orcId)
	if err != nil {
		log.Printf("Failed to convert OrcId to binary format. %s", err)
		return
	}

	binBibHash, err := bibHashToByteArray(bibHash)
	if err != nil {
		log.Printf("Failed to convert bibliographic hash to binary format. %s", err)
		return
	}

	transactOpts := bind.NewKeyedTransactor(st.PrivateKey)
	transaction, err := st.Contract.AddFeedback(transactOpts, binOrcId, binBibHash, relevance, presentation, methodology)

	if err != nil {
		log.Printf("Failed to deploy feedback transaction to ledger.")
	}

	return
}

// bibHashToByteArray converts a level 1 bibliographic hash (16 byte hex string) to its byte representation.
// Returns error if length of input string differes 32 or if it contains invalid characters.
func bibHashToByteArray(bibHash string) (data [16]byte, err error) {

	numChars := len(bibHash)

	if numChars != 32 {
		err = fmt.Errorf("Bibliographic hash string has wrong length. Should be 16 bytes in hex. Expected 32 characters but got %d.", numChars)
		return
	}

	decoded, err := hex.DecodeString(bibHash)
	if err != nil {
		err = fmt.Errorf("Bibliographic hash could not be decoded to binary representation. %s", err)
		return
	}

	copy(data[:], decoded)

	return
}

// orcIdToByteArray converts an OrcId string representation (e.g. "0000-0001-5000-000X") to a charachter encoded byte-array of length 16.
// Returns error if form of string does not follow standard 19-character form.
func orcIdToByteArray(orcId string) (data [16]byte, err error) {

	numChars := len(orcId)

	if numChars != 19 {
		err = fmt.Errorf("OrcId string has wrong length. Expected 19 characters but got %d.", numChars)
		return
	}

	// Local helper function that removes dashes from strings (as part of strings.Map).
	removeDash := func(r rune) rune {
		if r == '-' {
			return -1
		}
		return r
	}

	trimmedOrcId := strings.Map(removeDash, orcId)
	decoded := []byte(trimmedOrcId)
	copy(data[:], decoded)

	return
}
