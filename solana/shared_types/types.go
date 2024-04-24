// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package shared_types

import (
	solanago "github.com/blocto/solana-go-sdk/types"
	"github.com/coinbase/rosetta-sdk-go/types"
	bin "github.com/streamingfast/binary"
	"github.com/streamingfast/solana-go"
	dfuserpc "github.com/streamingfast/solana-go/rpc"
)

const (
	// NodeVersion is the version of geth we are using.
	NodeVersion = "1.4.17"

	// Blockchain is solanago.
	Blockchain string = "solana"

	// MainnetNetwork is the value of the network
	// in MainnetNetworkIdentifier.
	MainnetNetwork string = "mainnet"

	// TestnetNetwork is the value of the network
	// in TestnetNetworkIdentifier.
	TestnetNetwork string = "testnet"

	// DevnetNetwork is the value of the network
	// in DevnetNetworkIdentifier.
	DevnetNetwork string = "devnet"

	// Symbol is the symbol value
	// used in Currency.
	Symbol = "SOL"

	// Decimals is the decimals value
	// used in Currency.
	Decimals = 9

	// SuccessStatus is the status of any
	// Ethereum operation considered successful.
	SuccessStatus = "SUCCESS"

	// FailureStatus is the status of any
	// Ethereum operation considered unsuccessful.
	FailureStatus = "FAILURE"

	// HistoricalBalanceSupported is whether
	// historical balance is supported.
	HistoricalBalanceSupported = true

	// GenesisBlockIndex is the index of the
	// genesis block.
	GenesisBlockIndex = int64(0)

	Separator          = "__"
	WithNonceKey       = "with_nonce"
	PriorityFeeKey     = "priority_fee"
	SplSystemAccMapKey = "spl_system_acc_map"
	SplTokenAccMapKey  = "spl_token_acc_map"

	MainnetGenesisHash = "5eykt4UsFv8P8NJdTREpY1vzqKqZKvdpKuc147dw2N9d"
	TestnetGenesisHash = "4uhcVJyU9pJkvQyS88uRDiswHXSCkY3zQawwpjk2NsNY"
	DevnetGenesisHash  = "EtWTRABZaYq6iMfeYKouRu166VU2xqa1wcaWoxPkrZBG"
)

//op shared_types

const (
	System__Transfer                  = "System__Transfer"
	System__CreateAccount             = "System__CreateAccount"
	System__Assign                    = "System__Assign"
	System__CreateNonceAccount        = "System__CreateNonceAccount"
	System__AdvanceNonce              = "System__AdvanceNonce"
	System__WithdrawFromNonce         = "System__WithdrawFromNonce"
	System__AuthorizeNonce            = "System__AuthorizeNonce"
	System__Allocate                  = "System__Allocate"
	SplToken__Transfer                = "SplToken__Transfer"
	SplToken__InitializeMint          = "SplToken__InitializeMint"
	SplToken__InitializeAccount       = "SplToken__InitializeAccount"
	SplToken__CreateToken             = "SplToken__CreateToken"
	SplToken__CreateAccount           = "SplToken__CreateAccount"
	SplToken__Approve                 = "SplToken__Approve"
	SplToken__Revoke                  = "SplToken__Revoke"
	SplToken_MintTo                   = "SplToken_MintTo"
	SplToken_Burn                     = "SplToken_Burn"
	SplToken_CloseAccount             = "SplToken_CloseAccount"
	SplToken_FreezeAccount            = "SplToken_FreezeAccount"
	SplToken__TransferChecked         = "SplToken__TransferChecked"
	SplToken__TransferNew             = "SplToken__TransferNew"
	SplToken__TransferWithSystem      = "SplToken__TransferWithSystem"
	SplAssociatedTokenAccount__Create = "SplAssociatedTokenAccount__Create"
	Unknown                           = "Unknown"
	Stake__CreateStakeAccount         = "Stake__CreateStakeAccount"
	Stake__DelegateStake              = "Stake__DelegateStake"
	Stake__CreateStakeAndDelegate     = "Stake__CreateStakeAndDelegate"
	Stake__DeactivateStake            = "Stake__DeactivateStake"
	Stake__WithdrawStake              = "Stake__WithdrawStake"
	Stake__Merge                      = "Stake__Merge"
	Stake__Split                      = "Stake__Split"
	Stake__Authorize                  = "Stake__Authorize"
)

var (
	// MainnetGenesisBlockIdentifier is the *shared_types.BlockIdentifier
	// of the mainnet genesis block.
	MainnetGenesisBlockIdentifier = &types.BlockIdentifier{
		Hash:  MainnetGenesisHash,
		Index: GenesisBlockIndex,
	}

	// TestnetGenesisBlockIdentifier is the *shared_types.BlockIdentifier
	// of the testnet genesis block.
	TestnetGenesisBlockIdentifier = &types.BlockIdentifier{
		Hash:  TestnetGenesisHash,
		Index: GenesisBlockIndex,
	}
	// TestnetGenesisBlockIdentifier is the *shared_types.BlockIdentifier
	// of the testnet genesis block.
	DevnetGenesisHashBlockIdentifier = &types.BlockIdentifier{
		Hash:  DevnetGenesisHash,
		Index: GenesisBlockIndex,
	}

	// Currency is the *shared_types.Currency for all
	// Ethereum networks.
	Currency = &types.Currency{
		Symbol:   Symbol,
		Decimals: Decimals,
	}

	// OperationTypes are all suppoorted operation shared_types.
	OperationTypes = []string{
		System__Transfer,
		System__CreateAccount,
		System__Assign,
		System__CreateNonceAccount,
		System__AdvanceNonce,
		System__WithdrawFromNonce,
		System__AuthorizeNonce,
		System__Allocate,
		SplToken__Transfer,
		SplToken__InitializeMint,
		SplToken__InitializeAccount,
		SplToken__CreateToken,
		SplToken__CreateAccount,
		SplToken__Approve,
		SplToken__Revoke,
		SplToken_MintTo,
		SplToken_Burn,
		SplToken_CloseAccount,
		SplToken_FreezeAccount,
		SplToken__TransferChecked,
		SplToken__TransferNew,
		SplToken__TransferWithSystem,
		SplAssociatedTokenAccount__Create,
		Stake__CreateStakeAccount,
		Stake__DelegateStake,
		Stake__CreateStakeAndDelegate,
		Stake__DeactivateStake,
		Stake__WithdrawStake,
		Stake__Merge,
		Stake__Split,
		Stake__Authorize,
		Unknown,
	}

	// OperationStatuses are all supported operation statuses.
	OperationStatuses = []*types.OperationStatus{
		{
			Status:     SuccessStatus,
			Successful: true,
		},
		{
			Status:     FailureStatus,
			Successful: false,
		},
	}

	// CallMethods are all supported call methods.
	CallMethods = []string{
		"deregisterNode", "validatorExit", "getAccountInfo", "getBalance", "getBlockTime", "getClusterNodes", "getConfirmedBlock", "getConfirmedBlocks", "getConfirmedBlocksWithLimit", "getConfirmedSignaturesForAddress", "getConfirmedSignaturesForAddress2", "getConfirmedTransaction", "getEpochInfo", "getEpochSchedule", "getFeeCalculatorForBlockhash", "getFeeRateGovernor", "getFees", "getFirstAvailableBlock", "getGenesisHash", "getHealth", "getIdentity", "getInflationGovernor", "getInflationRate", "getLargestAccounts", "getLeaderSchedule", "getMinimumBalanceForRentExemption", "getMultipleAccounts", "getProgramAccounts", "getRecentBlockhash", "getSnapshotSlot", "getSignatureStatuses", "getSlot", "getSlotLeader", "getStorageTurn", "getStorageTurnRate", "getSlotsPerSegment", "getStoragePubkeysForSlot", "getSupply", "getTokenAccountBalance", "getTokenAccountsByDelegate", "getTokenAccountsByOwner", "getTokenSupply", "getTotalSupply", "getTransactionCount", "getVersion", "getVoteAccounts", "minimumLedgerSlot", "registerNode", "requestAirdrop", "sendTransaction", "simulateTransaction", "signVote",
	}
)

type TokenParsed struct {
	Decimals        uint64
	Amount          uint64
	MintAutority    solana.PublicKey
	FreezeAuthority solana.PublicKey
	AuthorityType   solana.PublicKey
	NewAuthority    solana.PublicKey
	M               byte
}

type ParsedInstructionMeta struct {
	Authority    string            `json:"authority,omitempty"`
	NewAuthority string            `json:"newAuthority,omitempty"`
	Source       string            `json:"source,omitempty"`
	Owner        string            `json:"owner,omitempty"`
	Account      string            `json:"account,omitempty"`
	Destination  string            `json:"destination,omitempty"`
	NewAccount   string            `json:"newAccount,omitempty"`
	Mint         string            `json:"mint,omitempty"`
	Decimals     uint8             `json:"decimals,omitempty"`
	TokenAmount  OpMetaTokenAmount `json:"tokenAmount,omitempty"`
	Amount       uint64            `json:"amount,omitempty"`
	Lamports     uint64            `json:"lamports,omitempty"`
	Space        uint64            `json:"space,omitempty"`
}
type OpMetaTokenAmount struct {
	Amount   string  `json:"amount,omitempty"`
	Decimals uint64  `json:"decimals,omitempty"`
	UiAmount float64 `json:"uiAmount,omitempty"`
}

type GetConfirmedBlockResult struct {
	Blockhash         solana.PublicKey             `json:"blockhash"`
	PreviousBlockhash solana.PublicKey             `json:"previousBlockhash"` // could be zeroes if ledger was clean-up and this is unavailable
	ParentSlot        bin.Uint64                   `json:"parentSlot"`
	Transactions      []dfuserpc.TransactionParsed `json:"transactions"`
	Rewards           []dfuserpc.BlockReward       `json:"rewards"`
	BlockTime         bin.Uint64                   `json:"blockTime,omitempty"`
}

type WithNonce struct {
	Account   string `json:"account"`
	Authority string `json:"authority,omitempty"`
}

type PriorityFee struct {
	MicroLamports string `json:"microLamports"`
}

type SplAccounts struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Mint        string `json:"mint"`
}

type Instruction struct {
	ProgramIDIndex uint64   `json:"programIdIndex"`
	Accounts       []uint64 `json:"accounts"`
	Data           string   `json:"data"`
}

type MessageHeader struct {
	NumRequiredSignatures       uint8 `json:"numRequiredSignatures"`
	NumReadonlySignedAccounts   uint8 `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts uint8 `json:"numReadonlyUnsignedAccounts"`
}

type TransactionMeta struct {
	Fee               uint64   `json:"fee"`
	PreBalances       []int64  `json:"preBalances"`
	PostBalances      []int64  `json:"postBalances"`
	LogMessages       []string `json:"logMesssages"`
	InnerInstructions []struct {
		Index        uint64        `json:"index"`
		Instructions []Instruction `json:"instructions"`
	} `json:"innerInstructions"`
	Err    interface{}            `json:"err"`
	Status map[string]interface{} `json:"status"`
}

type InstructionInfo struct {
	Info            map[string]interface{} `json:"info"`
	InstructionType string                 `json:"type"`
}

type Info struct {
	Authority     string              `json:"authority"`
	BlockHash     string              `json:"blockhash"`
	FeeCalculator FeeCalculatorString `json:"feeCalculator"`
}
type Parsed struct {
	Info Info `json:"info"`
}
type AccData struct {
	Parsed Parsed `json:"parsed"`
}

type FeeCalculator struct {
	LamportsPerSignature uint64 `json:"lamportsPerSignature"`
}

type FeeCalculatorString struct {
	LamportsPerSignature string `json:"lamportsPerSignature"`
}

type ParsedTransaction struct {
	Signatures []string      `json:"signatures"`
	Message    ParsedMessage `json:"message"`
}

type Reward struct {
	Pubkey      string `json:"pubkey"`
	Lamports    int64  `json:"lamports"`
	PostBalance uint64 `json:"postBalance"`
	RewardType  string `json:"rewardType"` // type of reward: "fee", "rent", "voting", "staking"
}

type GetConfirmBlockParsedResponse struct {
	Blockhash         string                      `json:"blockhash"`
	PreviousBlockhash string                      `json:"previousBlockhash"`
	ParentSlot        uint64                      `json:"parentSlot"`
	BlockTime         int64                       `json:"blockTime"`
	Transactions      []ParsedTransactionWithMeta `json:"transactions"`
	Rewards           []Reward                    `json:"rewards"`
}

type ParsedTransactionWithMeta struct {
	Meta        TransactionMeta   `json:"meta"`
	Transaction ParsedTransaction `json:"transaction"`
}

type Transaction struct {
	Signatures []string `json:"signatures"`
	Message    Message  `json:"message"`
}

type Message struct {
	Header          MessageHeader `json:"header"`
	AccountKeys     []string      `json:"accountKeys"`
	RecentBlockhash string        `json:"recentBlockhash"`
	Instructions    []Instruction `json:"instructions"`
}

type TokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       int32   `json:"decimals"`
	UIAmount       float64 `json:"uiAmount"`
	UIAmountString string  `json:"uiAmountString"`
}
type ParsedAccountInfo struct {
	Delegate        string      `json:"delegate"`
	DelegatedAmount TokenAmount `json:"delegatedAmount,omitempty"`
	IsInitialized   bool        `json:"isInitialized"`
	IsNative        bool        `json:"isNative"`
	Mint            string      `json:"mint"`
	Owner           string      `json:"owner"`
	TokenAmount     TokenAmount `json:"tokenAmount"`
}
type ParsedAccountData struct {
	AccountType string            `json:"accountType"`
	Info        ParsedAccountInfo `json:"info"`
}
type Data struct {
	Parsed  ParsedAccountData `json:"parsed"`
	Program string            `json:"program"`
}
type Account struct {
	Data       Data   `json:"data"`
	Executable bool   `json:"executable"`
	Lamports   int64  `json:"lamports"`
	Owner      string `json:"owner"`
	RentEpoch  int64  `json:"rentEpoch"`
}
type Accounts struct {
	Account Account `json:"account"`
	Pubkey  string  `json:"pubkey,omitempty"`
}

type ParsedMessage struct {
	Header          solanago.MessageHeader `json:"header"`
	AccountKeys     []ParsedAccKey         `json:"accountKeys"`
	RecentBlockhash string                 `json:"recentBlockhash"`
	Instructions    []ParsedInstruction    `json:"instructions"`
}

type ParsedInstruction struct {
	Accounts  []string         `json:"accounts,omitempty"`
	Data      string           `json:"data,omitempty"`
	Parsed    *InstructionInfo `json:"parsed,omitempty"`
	Program   string           `json:"program,omitempty"`
	ProgramID string           `json:"programId"`
}

type ParsedAccKey struct {
	PubKey     string `json:"pubkey"`
	IsSigner   bool   `json:"signer"`
	IsWritable bool   `json:"writable"`
}
