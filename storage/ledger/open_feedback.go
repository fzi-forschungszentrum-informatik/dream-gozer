// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ledger

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// OpenFeedbackABI is the input ABI used to generate the binding from.
const OpenFeedbackABI = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"_orcId\",\"type\":\"bytes16\"},{\"internalType\":\"bytes16\",\"name\":\"_bibHash\",\"type\":\"bytes16\"},{\"internalType\":\"uint8\",\"name\":\"_relevance\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"_presentation\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"_methodology\",\"type\":\"uint8\"}],\"name\":\"addFeedback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"_bibhash\",\"type\":\"bytes16\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getFeedbackByBibHash\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"serviceAddress_\",\"type\":\"address\"},{\"internalType\":\"bytes16\",\"name\":\"orcId_\",\"type\":\"bytes16\"},{\"internalType\":\"uint256\",\"name\":\"timestamp_\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"relevance_\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"presentation_\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"methodology_\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"_orcid\",\"type\":\"bytes16\"},{\"internalType\":\"uint256\",\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getFeedbackByOrcId\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"serviceAddress_\",\"type\":\"address\"},{\"internalType\":\"bytes16\",\"name\":\"bibHash_\",\"type\":\"bytes16\"},{\"internalType\":\"uint256\",\"name\":\"timestamp_\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"relevance_\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"presentation_\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"methodology_\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"_bibHash\",\"type\":\"bytes16\"}],\"name\":\"getFeedbackCountByBibHash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"count_\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"_orcid\",\"type\":\"bytes16\"}],\"name\":\"getFeedbackCountByOrcId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"count_\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"bytes16\",\"name\":\"_bibhash\",\"type\":\"bytes16\"}],\"name\":\"getTotalFeedbackByBibHash\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"feedbackCount_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"relevanceTotal_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"presentationTotal_\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"methodologyTotal_\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// OpenFeedbackFuncSigs maps the 4-byte function signature to its string representation.
var OpenFeedbackFuncSigs = map[string]string{
	"22b70c85": "addFeedback(bytes16,bytes16,uint8,uint8,uint8)",
	"aa80c06e": "getFeedbackByBibHash(bytes16,uint256)",
	"fd218468": "getFeedbackByOrcId(bytes16,uint256)",
	"3b6b3298": "getFeedbackCountByBibHash(bytes16)",
	"b8c2177c": "getFeedbackCountByOrcId(bytes16)",
	"3c68a762": "getTotalFeedbackByBibHash(bytes16)",
}

// OpenFeedbackBin is the compiled bytecode used for deploying new contracts.
var OpenFeedbackBin = "0x608060405234801561001057600080fd5b50600080546001600160a01b031916331790556108e8806100326000396000f3fe608060405234801561001057600080fd5b50600436106100625760003560e01c806322b70c85146100675780633b6b3298146100b55780633c68a762146100ee578063aa80c06e1461013b578063b8c2177c146101b3578063fd218468146101da575b600080fd5b6100b3600480360360a081101561007d57600080fd5b506001600160801b0319813581169160208101359091169060ff6040820135811691606081013582169160809091013516610207565b005b6100dc600480360360208110156100cb57600080fd5b50356001600160801b031916610531565b60408051918252519081900360200190f35b6101156004803603602081101561010457600080fd5b50356001600160801b03191661054d565b604080519485526020850193909352838301919091526060830152519081900360800190f35b6101686004803603604081101561015157600080fd5b506001600160801b0319813516906020013561066a565b604080516001600160a01b0390971687526001600160801b031990951660208701528585019390935260ff9182166060860152811660808501521660a0830152519081900360c00190f35b6100dc600480360360208110156101c957600080fd5b50356001600160801b031916610760565b610168600480360360408110156101f057600080fd5b506001600160801b0319813516906020013561077c565b61020f610877565b3381600001906001600160a01b031690816001600160a01b0316815250508581602001906001600160801b03191690816001600160801b031916815250508481604001906001600160801b03191690816001600160801b031916815250504281606001818152505083816080019060ff16908160ff1681525050828160a0019060ff16908160ff1681525050818160c0019060ff16908160ff168152505060026000866001600160801b0319166001600160801b0319168152602001908152602001600020819080600181540180825580915050906001820390600052602060002090600402016000909192909190915060008201518160000160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060208201518160010160006101000a8154816001600160801b03021916908360801c021790555060408201518160010160106101000a8154816001600160801b03021916908360801c02179055506060820151816002015560808201518160030160006101000a81548160ff021916908360ff16021790555060a08201518160030160016101000a81548160ff021916908360ff16021790555060c08201518160030160026101000a81548160ff021916908360ff16021790555050505060016000876001600160801b0319166001600160801b0319168152602001908152602001600020819080600181540180825580915050906001820390600052602060002090600402016000909192909190915060008201518160000160006101000a8154816001600160a01b0302191690836001600160a01b0316021790555060208201518160010160006101000a8154816001600160801b03021916908360801c021790555060408201518160010160106101000a8154816001600160801b03021916908360801c02179055506060820151816002015560808201518160030160006101000a81548160ff021916908360ff16021790555060a08201518160030160016101000a81548160ff021916908360ff16021790555060c08201518160030160026101000a81548160ff021916908360ff160217905550505050505050505050565b6001600160801b03191660009081526002602052604090205490565b6001600160801b03198116600090815260026020526040812054908080805b84811015610662576001600160801b03198616600090815260026020526040812080548390811061059957fe5b600091825260209091206003600490920201015460ff1611156105bd576001909301925b6001600160801b0319861660009081526002602052604081208054839081106105e257fe5b6000918252602090912060049091020160030154610100900460ff16111561060b576001909201915b6001600160801b03198616600090815260026020526040812080548390811061063057fe5b600091825260209091206004909102016003015462010000900460ff16111561065a576001909101905b60010161056c565b509193509193565b6001600160801b031982166000908152600260205260408120548190819081908190819087106106e1576040805162461bcd60e51b815260206004820152601f60248201527f53706563696669656420696e646578206973206f7574206f662072616e676500604482015290519081900360640190fd5b6001600160801b03198816600090815260026020526040812080548990811061070657fe5b600091825260209091206004909102018054600182015460028301546003909301546001600160a01b039092169c60809190911b9b5091995060ff8082169950610100820481169850620100009091041695509350505050565b6001600160801b03191660009081526001602052604090205490565b6001600160801b031982166000908152600160205260408120548190819081908190819087106107f3576040805162461bcd60e51b815260206004820152601f60248201527f53706563696669656420696e646578206973206f7574206f662072616e676500604482015290519081900360640190fd5b6001600160801b03198816600090815260016020526040812080548990811061081857fe5b600091825260209091206004909102018054600182015460028301546003909301546001600160a01b039092169c600160801b90910460801b9b5091995060ff8082169950610100820481169850620100009091041695509350505050565b6040805160e081018252600080825260208201819052918101829052606081018290526080810182905260a0810182905260c08101919091529056fea265627a7a723158209e2da53537a6825d679ef397c6113871fa9ac0b6ebae49a461b3d5b71a17eeef64736f6c634300050c0032"

// DeployOpenFeedback deploys a new Ethereum contract, binding an instance of OpenFeedback to it.
func DeployOpenFeedback(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *OpenFeedback, error) {
	parsed, err := abi.JSON(strings.NewReader(OpenFeedbackABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(OpenFeedbackBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &OpenFeedback{OpenFeedbackCaller: OpenFeedbackCaller{contract: contract}, OpenFeedbackTransactor: OpenFeedbackTransactor{contract: contract}, OpenFeedbackFilterer: OpenFeedbackFilterer{contract: contract}}, nil
}

// OpenFeedback is an auto generated Go binding around an Ethereum contract.
type OpenFeedback struct {
	OpenFeedbackCaller     // Read-only binding to the contract
	OpenFeedbackTransactor // Write-only binding to the contract
	OpenFeedbackFilterer   // Log filterer for contract events
}

// OpenFeedbackCaller is an auto generated read-only Go binding around an Ethereum contract.
type OpenFeedbackCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpenFeedbackTransactor is an auto generated write-only Go binding around an Ethereum contract.
type OpenFeedbackTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpenFeedbackFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type OpenFeedbackFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// OpenFeedbackSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type OpenFeedbackSession struct {
	Contract     *OpenFeedback     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// OpenFeedbackCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type OpenFeedbackCallerSession struct {
	Contract *OpenFeedbackCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// OpenFeedbackTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type OpenFeedbackTransactorSession struct {
	Contract     *OpenFeedbackTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// OpenFeedbackRaw is an auto generated low-level Go binding around an Ethereum contract.
type OpenFeedbackRaw struct {
	Contract *OpenFeedback // Generic contract binding to access the raw methods on
}

// OpenFeedbackCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type OpenFeedbackCallerRaw struct {
	Contract *OpenFeedbackCaller // Generic read-only contract binding to access the raw methods on
}

// OpenFeedbackTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type OpenFeedbackTransactorRaw struct {
	Contract *OpenFeedbackTransactor // Generic write-only contract binding to access the raw methods on
}

// NewOpenFeedback creates a new instance of OpenFeedback, bound to a specific deployed contract.
func NewOpenFeedback(address common.Address, backend bind.ContractBackend) (*OpenFeedback, error) {
	contract, err := bindOpenFeedback(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &OpenFeedback{OpenFeedbackCaller: OpenFeedbackCaller{contract: contract}, OpenFeedbackTransactor: OpenFeedbackTransactor{contract: contract}, OpenFeedbackFilterer: OpenFeedbackFilterer{contract: contract}}, nil
}

// NewOpenFeedbackCaller creates a new read-only instance of OpenFeedback, bound to a specific deployed contract.
func NewOpenFeedbackCaller(address common.Address, caller bind.ContractCaller) (*OpenFeedbackCaller, error) {
	contract, err := bindOpenFeedback(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &OpenFeedbackCaller{contract: contract}, nil
}

// NewOpenFeedbackTransactor creates a new write-only instance of OpenFeedback, bound to a specific deployed contract.
func NewOpenFeedbackTransactor(address common.Address, transactor bind.ContractTransactor) (*OpenFeedbackTransactor, error) {
	contract, err := bindOpenFeedback(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &OpenFeedbackTransactor{contract: contract}, nil
}

// NewOpenFeedbackFilterer creates a new log filterer instance of OpenFeedback, bound to a specific deployed contract.
func NewOpenFeedbackFilterer(address common.Address, filterer bind.ContractFilterer) (*OpenFeedbackFilterer, error) {
	contract, err := bindOpenFeedback(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &OpenFeedbackFilterer{contract: contract}, nil
}

// bindOpenFeedback binds a generic wrapper to an already deployed contract.
func bindOpenFeedback(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(OpenFeedbackABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OpenFeedback *OpenFeedbackRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OpenFeedback.Contract.OpenFeedbackCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OpenFeedback *OpenFeedbackRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpenFeedback.Contract.OpenFeedbackTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OpenFeedback *OpenFeedbackRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OpenFeedback.Contract.OpenFeedbackTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_OpenFeedback *OpenFeedbackCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _OpenFeedback.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_OpenFeedback *OpenFeedbackTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _OpenFeedback.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_OpenFeedback *OpenFeedbackTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _OpenFeedback.Contract.contract.Transact(opts, method, params...)
}

// GetFeedbackByBibHash is a free data retrieval call binding the contract method 0xaa80c06e.
//
// Solidity: function getFeedbackByBibHash(bytes16 _bibhash, uint256 _index) constant returns(address serviceAddress_, bytes16 orcId_, uint256 timestamp_, uint8 relevance_, uint8 presentation_, uint8 methodology_)
func (_OpenFeedback *OpenFeedbackCaller) GetFeedbackByBibHash(opts *bind.CallOpts, _bibhash [16]byte, _index *big.Int) (struct {
	ServiceAddress common.Address
	OrcId          [16]byte
	Timestamp      *big.Int
	Relevance      uint8
	Presentation   uint8
	Methodology    uint8
}, error) {
	ret := new(struct {
		ServiceAddress common.Address
		OrcId          [16]byte
		Timestamp      *big.Int
		Relevance      uint8
		Presentation   uint8
		Methodology    uint8
	})
	out := ret
	err := _OpenFeedback.contract.Call(opts, out, "getFeedbackByBibHash", _bibhash, _index)
	return *ret, err
}

// GetFeedbackByBibHash is a free data retrieval call binding the contract method 0xaa80c06e.
//
// Solidity: function getFeedbackByBibHash(bytes16 _bibhash, uint256 _index) constant returns(address serviceAddress_, bytes16 orcId_, uint256 timestamp_, uint8 relevance_, uint8 presentation_, uint8 methodology_)
func (_OpenFeedback *OpenFeedbackSession) GetFeedbackByBibHash(_bibhash [16]byte, _index *big.Int) (struct {
	ServiceAddress common.Address
	OrcId          [16]byte
	Timestamp      *big.Int
	Relevance      uint8
	Presentation   uint8
	Methodology    uint8
}, error) {
	return _OpenFeedback.Contract.GetFeedbackByBibHash(&_OpenFeedback.CallOpts, _bibhash, _index)
}

// GetFeedbackByBibHash is a free data retrieval call binding the contract method 0xaa80c06e.
//
// Solidity: function getFeedbackByBibHash(bytes16 _bibhash, uint256 _index) constant returns(address serviceAddress_, bytes16 orcId_, uint256 timestamp_, uint8 relevance_, uint8 presentation_, uint8 methodology_)
func (_OpenFeedback *OpenFeedbackCallerSession) GetFeedbackByBibHash(_bibhash [16]byte, _index *big.Int) (struct {
	ServiceAddress common.Address
	OrcId          [16]byte
	Timestamp      *big.Int
	Relevance      uint8
	Presentation   uint8
	Methodology    uint8
}, error) {
	return _OpenFeedback.Contract.GetFeedbackByBibHash(&_OpenFeedback.CallOpts, _bibhash, _index)
}

// GetFeedbackByOrcId is a free data retrieval call binding the contract method 0xfd218468.
//
// Solidity: function getFeedbackByOrcId(bytes16 _orcid, uint256 _index) constant returns(address serviceAddress_, bytes16 bibHash_, uint256 timestamp_, uint8 relevance_, uint8 presentation_, uint8 methodology_)
func (_OpenFeedback *OpenFeedbackCaller) GetFeedbackByOrcId(opts *bind.CallOpts, _orcid [16]byte, _index *big.Int) (struct {
	ServiceAddress common.Address
	BibHash        [16]byte
	Timestamp      *big.Int
	Relevance      uint8
	Presentation   uint8
	Methodology    uint8
}, error) {
	ret := new(struct {
		ServiceAddress common.Address
		BibHash        [16]byte
		Timestamp      *big.Int
		Relevance      uint8
		Presentation   uint8
		Methodology    uint8
	})
	out := ret
	err := _OpenFeedback.contract.Call(opts, out, "getFeedbackByOrcId", _orcid, _index)
	return *ret, err
}

// GetFeedbackByOrcId is a free data retrieval call binding the contract method 0xfd218468.
//
// Solidity: function getFeedbackByOrcId(bytes16 _orcid, uint256 _index) constant returns(address serviceAddress_, bytes16 bibHash_, uint256 timestamp_, uint8 relevance_, uint8 presentation_, uint8 methodology_)
func (_OpenFeedback *OpenFeedbackSession) GetFeedbackByOrcId(_orcid [16]byte, _index *big.Int) (struct {
	ServiceAddress common.Address
	BibHash        [16]byte
	Timestamp      *big.Int
	Relevance      uint8
	Presentation   uint8
	Methodology    uint8
}, error) {
	return _OpenFeedback.Contract.GetFeedbackByOrcId(&_OpenFeedback.CallOpts, _orcid, _index)
}

// GetFeedbackByOrcId is a free data retrieval call binding the contract method 0xfd218468.
//
// Solidity: function getFeedbackByOrcId(bytes16 _orcid, uint256 _index) constant returns(address serviceAddress_, bytes16 bibHash_, uint256 timestamp_, uint8 relevance_, uint8 presentation_, uint8 methodology_)
func (_OpenFeedback *OpenFeedbackCallerSession) GetFeedbackByOrcId(_orcid [16]byte, _index *big.Int) (struct {
	ServiceAddress common.Address
	BibHash        [16]byte
	Timestamp      *big.Int
	Relevance      uint8
	Presentation   uint8
	Methodology    uint8
}, error) {
	return _OpenFeedback.Contract.GetFeedbackByOrcId(&_OpenFeedback.CallOpts, _orcid, _index)
}

// GetFeedbackCountByBibHash is a free data retrieval call binding the contract method 0x3b6b3298.
//
// Solidity: function getFeedbackCountByBibHash(bytes16 _bibHash) constant returns(uint256 count_)
func (_OpenFeedback *OpenFeedbackCaller) GetFeedbackCountByBibHash(opts *bind.CallOpts, _bibHash [16]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OpenFeedback.contract.Call(opts, out, "getFeedbackCountByBibHash", _bibHash)
	return *ret0, err
}

// GetFeedbackCountByBibHash is a free data retrieval call binding the contract method 0x3b6b3298.
//
// Solidity: function getFeedbackCountByBibHash(bytes16 _bibHash) constant returns(uint256 count_)
func (_OpenFeedback *OpenFeedbackSession) GetFeedbackCountByBibHash(_bibHash [16]byte) (*big.Int, error) {
	return _OpenFeedback.Contract.GetFeedbackCountByBibHash(&_OpenFeedback.CallOpts, _bibHash)
}

// GetFeedbackCountByBibHash is a free data retrieval call binding the contract method 0x3b6b3298.
//
// Solidity: function getFeedbackCountByBibHash(bytes16 _bibHash) constant returns(uint256 count_)
func (_OpenFeedback *OpenFeedbackCallerSession) GetFeedbackCountByBibHash(_bibHash [16]byte) (*big.Int, error) {
	return _OpenFeedback.Contract.GetFeedbackCountByBibHash(&_OpenFeedback.CallOpts, _bibHash)
}

// GetFeedbackCountByOrcId is a free data retrieval call binding the contract method 0xb8c2177c.
//
// Solidity: function getFeedbackCountByOrcId(bytes16 _orcid) constant returns(uint256 count_)
func (_OpenFeedback *OpenFeedbackCaller) GetFeedbackCountByOrcId(opts *bind.CallOpts, _orcid [16]byte) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _OpenFeedback.contract.Call(opts, out, "getFeedbackCountByOrcId", _orcid)
	return *ret0, err
}

// GetFeedbackCountByOrcId is a free data retrieval call binding the contract method 0xb8c2177c.
//
// Solidity: function getFeedbackCountByOrcId(bytes16 _orcid) constant returns(uint256 count_)
func (_OpenFeedback *OpenFeedbackSession) GetFeedbackCountByOrcId(_orcid [16]byte) (*big.Int, error) {
	return _OpenFeedback.Contract.GetFeedbackCountByOrcId(&_OpenFeedback.CallOpts, _orcid)
}

// GetFeedbackCountByOrcId is a free data retrieval call binding the contract method 0xb8c2177c.
//
// Solidity: function getFeedbackCountByOrcId(bytes16 _orcid) constant returns(uint256 count_)
func (_OpenFeedback *OpenFeedbackCallerSession) GetFeedbackCountByOrcId(_orcid [16]byte) (*big.Int, error) {
	return _OpenFeedback.Contract.GetFeedbackCountByOrcId(&_OpenFeedback.CallOpts, _orcid)
}

// GetTotalFeedbackByBibHash is a free data retrieval call binding the contract method 0x3c68a762.
//
// Solidity: function getTotalFeedbackByBibHash(bytes16 _bibhash) constant returns(uint256 feedbackCount_, uint256 relevanceTotal_, uint256 presentationTotal_, uint256 methodologyTotal_)
func (_OpenFeedback *OpenFeedbackCaller) GetTotalFeedbackByBibHash(opts *bind.CallOpts, _bibhash [16]byte) (struct {
	FeedbackCount     *big.Int
	RelevanceTotal    *big.Int
	PresentationTotal *big.Int
	MethodologyTotal  *big.Int
}, error) {
	ret := new(struct {
		FeedbackCount     *big.Int
		RelevanceTotal    *big.Int
		PresentationTotal *big.Int
		MethodologyTotal  *big.Int
	})
	out := ret
	err := _OpenFeedback.contract.Call(opts, out, "getTotalFeedbackByBibHash", _bibhash)
	return *ret, err
}

// GetTotalFeedbackByBibHash is a free data retrieval call binding the contract method 0x3c68a762.
//
// Solidity: function getTotalFeedbackByBibHash(bytes16 _bibhash) constant returns(uint256 feedbackCount_, uint256 relevanceTotal_, uint256 presentationTotal_, uint256 methodologyTotal_)
func (_OpenFeedback *OpenFeedbackSession) GetTotalFeedbackByBibHash(_bibhash [16]byte) (struct {
	FeedbackCount     *big.Int
	RelevanceTotal    *big.Int
	PresentationTotal *big.Int
	MethodologyTotal  *big.Int
}, error) {
	return _OpenFeedback.Contract.GetTotalFeedbackByBibHash(&_OpenFeedback.CallOpts, _bibhash)
}

// GetTotalFeedbackByBibHash is a free data retrieval call binding the contract method 0x3c68a762.
//
// Solidity: function getTotalFeedbackByBibHash(bytes16 _bibhash) constant returns(uint256 feedbackCount_, uint256 relevanceTotal_, uint256 presentationTotal_, uint256 methodologyTotal_)
func (_OpenFeedback *OpenFeedbackCallerSession) GetTotalFeedbackByBibHash(_bibhash [16]byte) (struct {
	FeedbackCount     *big.Int
	RelevanceTotal    *big.Int
	PresentationTotal *big.Int
	MethodologyTotal  *big.Int
}, error) {
	return _OpenFeedback.Contract.GetTotalFeedbackByBibHash(&_OpenFeedback.CallOpts, _bibhash)
}

// AddFeedback is a paid mutator transaction binding the contract method 0x22b70c85.
//
// Solidity: function addFeedback(bytes16 _orcId, bytes16 _bibHash, uint8 _relevance, uint8 _presentation, uint8 _methodology) returns()
func (_OpenFeedback *OpenFeedbackTransactor) AddFeedback(opts *bind.TransactOpts, _orcId [16]byte, _bibHash [16]byte, _relevance uint8, _presentation uint8, _methodology uint8) (*types.Transaction, error) {
	return _OpenFeedback.contract.Transact(opts, "addFeedback", _orcId, _bibHash, _relevance, _presentation, _methodology)
}

// AddFeedback is a paid mutator transaction binding the contract method 0x22b70c85.
//
// Solidity: function addFeedback(bytes16 _orcId, bytes16 _bibHash, uint8 _relevance, uint8 _presentation, uint8 _methodology) returns()
func (_OpenFeedback *OpenFeedbackSession) AddFeedback(_orcId [16]byte, _bibHash [16]byte, _relevance uint8, _presentation uint8, _methodology uint8) (*types.Transaction, error) {
	return _OpenFeedback.Contract.AddFeedback(&_OpenFeedback.TransactOpts, _orcId, _bibHash, _relevance, _presentation, _methodology)
}

// AddFeedback is a paid mutator transaction binding the contract method 0x22b70c85.
//
// Solidity: function addFeedback(bytes16 _orcId, bytes16 _bibHash, uint8 _relevance, uint8 _presentation, uint8 _methodology) returns()
func (_OpenFeedback *OpenFeedbackTransactorSession) AddFeedback(_orcId [16]byte, _bibHash [16]byte, _relevance uint8, _presentation uint8, _methodology uint8) (*types.Transaction, error) {
	return _OpenFeedback.Contract.AddFeedback(&_OpenFeedback.TransactOpts, _orcId, _bibHash, _relevance, _presentation, _methodology)
}
