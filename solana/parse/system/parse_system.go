package system

import (
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/ghostiam/binstruct"
	stypes "github.com/imerkle/rosetta-solana-go/solana/shared_types"
	"log"
)

type Instruction uint32

const (
	InstructionCreateAccount Instruction = iota
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

func ParseSystem(ins types.Instruction) (stypes.ParsedInstruction, error) {
	var parsedInstruction stypes.ParsedInstruction
	var err error
	var s struct {
		Instruction Instruction
	}
	err = binstruct.UnmarshalLE(ins.Data, &s)
	var instructionType string
	var parsedInfo map[string]interface{}
	switch s.Instruction {
	case InstructionCreateAccount:
		var a CreateAccountInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionCreateAccount: %v", err)
		}
		instructionType = "createAccount"
		parsedInfo = map[string]interface{}{
			"source":     ins.Accounts[0].PubKey.ToBase58(),
			"newAccount": ins.Accounts[1].PubKey.ToBase58(),
			"lamports":   a.Lamports,
			"space":      a.Space,
			"owner":      a.Owner.ToBase58(),
		}
		break
	case InstructionAssign:
		var a AssignInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionAssign: %v", err)
		}
		instructionType = "assign"
		parsedInfo = map[string]interface{}{
			"account": ins.Accounts[0].PubKey.ToBase58(),
			"owner":   a.AssignToProgramID.ToBase58(),
		}
		break
	case InstructionTransfer:
		var a TransferInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionTransfer: %v", err)
		}
		instructionType = "transfer"
		parsedInfo = map[string]interface{}{
			"source":      ins.Accounts[0].PubKey.ToBase58(),
			"destination": ins.Accounts[1].PubKey.ToBase58(),
			"lamports":    a.Lamports,
		}
		parsedInstruction.Parsed = &stypes.InstructionInfo{}
		break
	case InstructionCreateAccountWithSeed:
		var a CreateAccountWithSeedInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionCreateAccountWithSeed: %v", err)
		}
		instructionType = "createAccountWithSeed"
		parsedInfo = map[string]interface{}{
			"source":     ins.Accounts[0].PubKey.ToBase58(),
			"newAccount": ins.Accounts[1].PubKey.ToBase58(),
			"base":       a.Base,
			"seed":       a.Seed,
			"space":      a.Space,
			"lamports":   a.Lamports,
			"owner":      a.ProgramID,
		}
		break
	case InstructionAdvanceNonceAccount:
		var a AdvanceNonceAccountInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionAdvanceNonceAccount: %v", err)
		}
		instructionType = "advanceNonce"
		parsedInfo = map[string]interface{}{
			"nonceAccount":            ins.Accounts[0].PubKey.ToBase58(),
			"recentBlockhashesSysvar": ins.Accounts[1].PubKey.ToBase58(),
			"nonceAuthority":          ins.Accounts[2].PubKey.ToBase58(),
		}
		break
	case InstructionWithdrawNonceAccount:
		var a WithdrawNonceAccountInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionWithdrawNonceAccount: %v", err)
		}
		instructionType = "withdrawFromNonce"
		parsedInfo = map[string]interface{}{
			"nonceAccount":            ins.Accounts[0].PubKey.ToBase58(),
			"destination":             ins.Accounts[1].PubKey.ToBase58(),
			"recentBlockhashesSysvar": ins.Accounts[2].PubKey.ToBase58(),
			"rentSysvar":              ins.Accounts[3].PubKey.ToBase58(),
			"nonceAuthority":          ins.Accounts[4].PubKey.ToBase58(),
			"lamports":                a.Lamports,
		}
		break
	case InstructionInitializeNonceAccount:
		var a InitializeNonceAccountInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionInitializeNonceAccount: %v", err)
		}
		instructionType = "initializeNonce"
		parsedInfo = map[string]interface{}{
			"nonceAccount":            ins.Accounts[0].PubKey.ToBase58(),
			"recentBlockhashesSysvar": ins.Accounts[1].PubKey.ToBase58(),
			"rentSysvar":              ins.Accounts[2].PubKey.ToBase58(),
			"nonceAuthority":          a.Auth.ToBase58(),
		}
		break
	case InstructionAuthorizeNonceAccount:
		var a AuthorizeNonceAccountInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionAuthorizeNonceAccount: %v", err)
		}
		instructionType = "authorizeNonce"
		parsedInfo = map[string]interface{}{
			"nonceAccount":   ins.Accounts[0].PubKey.ToBase58(),
			"nonceAuthority": ins.Accounts[1].PubKey.ToBase58(),
			"rentSysvar":     ins.Accounts[2].PubKey.ToBase58(),
			"newAuthorized":  a.Auth.ToBase58(),
		}
		break
	case InstructionAllocate:
		var a AllocateInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionAllocate: %v", err)
		}
		instructionType = "allocate"
		parsedInfo = map[string]interface{}{
			"account": ins.Accounts[0].PubKey.ToBase58(),
			"space":   a.Space,
		}
		break
	case InstructionAllocateWithSeed:
		var a AllocateWithSeedInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionAllocateWithSeed: %v", err)
		}
		instructionType = "allocateWithSeed"
		parsedInfo = map[string]interface{}{
			"account": ins.Accounts[0].PubKey.ToBase58(),
			"base":    a.Base,
			"seed":    a.Seed,
			"space":   a.Space,
			"owner":   a.ProgramID.ToBase58(),
		}
		break
	case InstructionAssignWithSeed:
		var a AssignWithSeedInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionAssignWithSeed: %v", err)
		}
		instructionType = "assignWithSeed"
		parsedInfo = map[string]interface{}{
			"account": ins.Accounts[0].PubKey.ToBase58(),
			"base":    a.Base,
			"seed":    a.Seed,
			"owner":   a.AssignToProgramID.ToBase58(),
		}
		break
	case InstructionTransferWithSeed:
		var a TransferWithSeedInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionTransferWithSeed: %v", err)
		}
		instructionType = "transferWithSeed"
		parsedInfo = map[string]interface{}{
			"source":      ins.Accounts[0].PubKey.ToBase58(),
			"sourceBase":  ins.Accounts[1].PubKey.ToBase58(),
			"destination": ins.Accounts[2].PubKey.ToBase58(),
			"lamports":    a.Lamports,
			"sourceSeed":  a.Seed,
			"sourceOwner": a.ProgramID.ToBase58(),
		}
		break
	}
	parsedInstruction.Parsed = &stypes.InstructionInfo{
		Info:            parsedInfo,
		InstructionType: instructionType,
	}
	return parsedInstruction, err
}

type AssignInstruction struct {
	Instruction       Instruction
	AssignToProgramID common.PublicKey
}

type TransferInstruction struct {
	Instruction Instruction
	Lamports    uint64
}

type CreateAccountWithSeedInstruction struct {
	Instruction Instruction
	Base        common.PublicKey
	SeedLen     uint64
	Seed        string
	Lamports    uint64
	Space       uint64
	ProgramID   common.PublicKey
}

type AdvanceNonceAccountInstruction struct {
	Instruction Instruction
}
type WithdrawNonceAccountInstruction struct {
	Instruction Instruction
	Lamports    uint64
}

type InitializeNonceAccountInstruction struct {
	Instruction Instruction
	Auth        common.PublicKey
}
type AuthorizeNonceAccountInstruction struct {
	Instruction Instruction
	Auth        common.PublicKey
}
type AllocateInstruction struct {
	Instruction Instruction
	Space       uint64
}
type AllocateWithSeedInstruction struct {
	Instruction Instruction
	Base        common.PublicKey
	SeedLen     uint64
	Seed        string
	Space       uint64
	ProgramID   common.PublicKey
}
type AssignWithSeedInstruction struct {
	Instruction       Instruction
	Base              common.PublicKey
	SeedLen           uint64
	Seed              string
	AssignToProgramID common.PublicKey
}
type TransferWithSeedInstruction struct {
	Instruction Instruction
	Lamports    uint64
	SeedLen     uint64
	Seed        string
	ProgramID   common.PublicKey
}

type CreateAccountInstruction struct {
	Instruction Instruction
	Lamports    uint64
	Space       uint64
	Owner       common.PublicKey
}
