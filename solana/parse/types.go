package parse

import (
	"github.com/blocto/solana-go-sdk/common"
	stypes "github.com/imerkle/rosetta-solana-go/solana/shared_types"
)

type InitializeMintInstruction struct {
	Instruction     stypes.Instruction
	Decimals        uint8
	MintAuthority   common.PublicKey
	Option          bool
	FreezeAuthority common.PublicKey
}
type InitializeAccountInstruction struct {
	Instruction stypes.Instruction
}
type InitializeMultisigInstruction struct {
	Instruction     stypes.Instruction
	MinimumRequired uint8
}
type ApproveInstruction struct {
	Instruction stypes.Instruction
	Amount      uint64
}
type RevokeInstruction struct {
	Instruction stypes.Instruction
}
type MintToInstruction struct {
	Instruction stypes.Instruction
	Amount      uint64
}
type BurnInstruction struct {
	Instruction stypes.Instruction
	Amount      uint64
}
type CloseAccountInstruction struct {
	Instruction stypes.Instruction
}
type FreezeAccountInstruction struct {
	Instruction stypes.Instruction
}
type ThawAccountInstruction struct {
	Instruction stypes.Instruction
}
type TransferCheckedInstruction struct {
	Instruction InstructionToken
	Amount      uint64
	Decimals    uint8
}
type ApproveCheckedInstruction struct {
	Instruction stypes.Instruction
	Amount      uint64
	Decimals    uint8
}
type MintToCheckedInstruction struct {
	Instruction stypes.Instruction
	Amount      uint64
	Decimals    uint8
}
type BurnCheckedInstruction struct {
	Instruction stypes.Instruction
	Amount      uint64
	Decimals    uint8
}
type InitializeAccount2Instruction struct {
	Instruction stypes.Instruction
	Owner       common.PublicKey
}

type UiTokenAmount struct {
	UiAmount float64
	Decimals uint8
	Amount   string
}

type TokenTransferInstruction struct {
	Instruction InstructionToken
	Amount      uint64
}

type InstructionInt uint32

const (
	InstructionCreateAccount InstructionInt = iota
	InstructionAssign
	InstructionTransfer
	InstructionCreateAccountWithSeed
	InstructionAdvanceNonceAccount
	InstructionWithdrawNonceAccount
	InstructionInitializeNonceAccount
	InstructionAuthorizeNonceAccount
	InstructionAllocate
	InstructionAllocateWithSeed
	InstructionAssignWithSeed
	InstructionTransferWithSeed
	InstructionUpgradeNonceAccount
)

type AssignInstruction struct {
	Instruction       stypes.Instruction
	AssignToProgramID common.PublicKey
}

type TransferInstruction struct {
	Instruction InstructionInt
	Lamports    uint64
}

type CreateAccountWithSeedInstruction struct {
	Instruction stypes.Instruction
	Base        common.PublicKey
	SeedLen     uint64
	Seed        string
	Lamports    uint64
	Space       uint64
	ProgramID   common.PublicKey
}

type AdvanceNonceAccountInstruction struct {
	Instruction stypes.Instruction
}
type WithdrawNonceAccountInstruction struct {
	Instruction stypes.Instruction
	Lamports    uint64
}

type InitializeNonceAccountInstruction struct {
	Instruction stypes.Instruction
	Auth        common.PublicKey
}
type AuthorizeNonceAccountInstruction struct {
	Instruction stypes.Instruction
	Auth        common.PublicKey
}
type AllocateInstruction struct {
	Instruction stypes.Instruction
	Space       uint64
}
type AllocateWithSeedInstruction struct {
	Instruction stypes.Instruction
	Base        common.PublicKey
	SeedLen     uint64
	Seed        string
	Space       uint64
	ProgramID   common.PublicKey
}
type AssignWithSeedInstruction struct {
	Instruction       stypes.Instruction
	Base              common.PublicKey
	SeedLen           uint64
	Seed              string
	AssignToProgramID common.PublicKey
}
type TransferWithSeedInstruction struct {
	Instruction stypes.Instruction
	Lamports    uint64
	SeedLen     uint64
	Seed        string
	ProgramID   common.PublicKey
}

type CreateAccountInstruction struct {
	Instruction stypes.Instruction
	Lamports    uint64
	Space       uint64
	Owner       common.PublicKey
}
