package solanago

import (
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/types"
	"github.com/ghostiam/binstruct"
	"log"
	"math"
)

type InstructionToken uint8

const (
	InstructionInitializeMint InstructionToken = iota
	InstructionInitializeAccount
	InstructionInitializeMultisig
	InstructionTransferToken
	InstructionApprove
	InstructionRevoke
	InstructionSetAuthority
	InstructionMintTo
	InstructionBurn
	InstructionCloseAccount
	InstructionFreezeAccount
	InstructionThawAccount
	InstructionTransferChecked
	InstructionApproveChecked
	InstructionMintToChecked
	InstructionBurnChecked
	InstructionInitializeAccount2
)

type InitializeMintInstruction struct {
	Instruction     Instruction
	Decimals        uint8
	MintAuthority   common.PublicKey
	Option          bool
	FreezeAuthority common.PublicKey
}
type InitializeAccountInstruction struct {
	Instruction Instruction
}
type InitializeMultisigInstruction struct {
	Instruction     Instruction
	MinimumRequired uint8
}
type ApproveInstruction struct {
	Instruction Instruction
	Amount      uint64
}
type RevokeInstruction struct {
	Instruction Instruction
}
type MintToInstruction struct {
	Instruction Instruction
	Amount      uint64
}
type BurnInstruction struct {
	Instruction Instruction
	Amount      uint64
}
type CloseAccountInstruction struct {
	Instruction Instruction
}
type FreezeAccountInstruction struct {
	Instruction Instruction
}
type ThawAccountInstruction struct {
	Instruction Instruction
}
type TransferCheckedInstruction struct {
	Instruction InstructionToken
	Amount      uint64
	Decimals    uint8
}
type ApproveCheckedInstruction struct {
	Instruction Instruction
	Amount      uint64
	Decimals    uint8
}
type MintToCheckedInstruction struct {
	Instruction Instruction
	Amount      uint64
	Decimals    uint8
}
type BurnCheckedInstruction struct {
	Instruction Instruction
	Amount      uint64
	Decimals    uint8
}
type InitializeAccount2Instruction struct {
	Instruction Instruction
	Owner       common.PublicKey
}

func ParseToken(ins types.Instruction) (ParsedInstruction, error) {
	var parsedInstruction ParsedInstruction
	var err error
	var s struct {
		Instruction InstructionToken
	}
	err = binstruct.UnmarshalLE(ins.Data, &s)
	var instructionType string
	var parsedInfo map[string]interface{}
	switch s.Instruction {
	case InstructionInitializeMint:
		var a InitializeMintInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "initializeMint"
		parsedInfo = map[string]interface{}{
			"mint":            ins.Accounts[0].PubKey.ToBase58(),
			"decimals":        a.Decimals,
			"mintAuthority":   a.MintAuthority.ToBase58(),
			"rentSysvar":      ins.Accounts[1].PubKey.ToBase58(),
			"freezeAuthority": a.FreezeAuthority.ToBase58(),
		}
		break
	case InstructionInitializeAccount:
		var a InitializeAccountInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "initializeAccount"
		parsedInfo = map[string]interface{}{
			"account":    ins.Accounts[0].PubKey.ToBase58(),
			"mint":       ins.Accounts[1].PubKey.ToBase58(),
			"owner":      ins.Accounts[2].PubKey.ToBase58(),
			"rentSysvar": ins.Accounts[3].PubKey.ToBase58(),
		}
		break
	//case InstructionInitializeMultisig:
	//	var a solanago.TransferInstruction
	//	err = binstruct.UnmarshalLE(ins.Data, &a)
	//	instructionType = "initializeMultisig"
	//	var signers []string
	//	for _, v := range ins.Accounts[2:] {
	//		signers = append(signers, v.PubKey.ToBase58())
	//	}
	//	parsedInfo = map[string]interface{}{
	//		"multisig":   ins.Accounts[0].PubKey.ToBase58(),
	//		"rentSysvar": ins.Accounts[1].PubKey.ToBase58(),
	//		"signers":    signers,
	//		"m":          a.Amount,
	//	}
	//	break
	case InstructionTransferToken:
		var a TokenTransferInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionTransferToken: %v", err)
		}
		instructionType = "transfer"
		parsedInfo = map[string]interface{}{
			"source":      ins.Accounts[0].PubKey.ToBase58(),
			"destination": ins.Accounts[1].PubKey.ToBase58(),
			"amount":      a.Amount,
		}
		parsedInfo = parse_signers(parsedInfo, 2, ins.Accounts, "authority", "multisigAuthority")
		break
	//case InstructionApprove:
	//	var a ApproveInstruction
	//	err = binstruct.UnmarshalLE(ins.Data, &a)
	//	instructionType = "approve"
	//	parsedInfo = map[string]interface{}{
	//		"source":   ins.Accounts[0].PubKey.ToBase58(),
	//		"delegate": ins.Accounts[1].PubKey.ToBase58(),
	//		"amount":   a.Amount,
	//	}
	//	parsedInfo = parse_signers(parsedInfo, 2, ins.Accounts, "owner", "multisigOwner")
	//	break
	//case InstructionRevoke:
	//	var a RevokeInstruction
	//	err = binstruct.UnmarshalLE(ins.Data, &a)
	//	instructionType = "revoke"
	//	parsedInfo = map[string]interface{}{
	//		"source": ins.Accounts[0].PubKey.ToBase58(),
	//	}
	//	parsedInfo = parse_signers(parsedInfo, 1, ins.Accounts, "owner", "multisigOwner")
	//
	//	break
	////case InstructionSetAuthority:
	////	break
	//case InstructionMintTo:
	//	var a MintToInstruction
	//	err = binstruct.UnmarshalLE(ins.Data, &a)
	//	instructionType = "mintTo"
	//	parsedInfo = map[string]interface{}{
	//		"mint":    ins.Accounts[0].PubKey.ToBase58(),
	//		"account": ins.Accounts[1].PubKey.ToBase58(),
	//		"amount":  a.Amount,
	//	}
	//	parsedInfo = parse_signers(parsedInfo, 2, ins.Accounts, "mintAuthority", "multisigMintAuthority")
	//
	//	break
	//case InstructionCloseAccount:
	//	var a CloseAccountInstruction
	//	err = binstruct.UnmarshalLE(ins.Data, &a)
	//	instructionType = "closeAccount"
	//	parsedInfo = map[string]interface{}{
	//		"account":     ins.Accounts[0].PubKey.ToBase58(),
	//		"destination": ins.Accounts[1].PubKey.ToBase58(),
	//	}
	//	parsedInfo = parse_signers(parsedInfo, 2, ins.Accounts, "owner", "multisigOwner")
	//
	//	break
	//case InstructionFreezeAccount:
	//	var a FreezeAccountInstruction
	//	err = binstruct.UnmarshalLE(ins.Data, &a)
	//	instructionType = "freezeAccount"
	//	parsedInfo = map[string]interface{}{
	//		"account": ins.Accounts[0].PubKey.ToBase58(),
	//		"mint":    ins.Accounts[1].PubKey.ToBase58(),
	//	}
	//	parsedInfo = parse_signers(parsedInfo, 2, ins.Accounts, "freezeAuthority", "multisigFreezeAuthority")
	//
	//	break
	//case InstructionThawAccount:
	//	var a ThawAccountInstruction
	//	err = binstruct.UnmarshalLE(ins.Data, &a)
	//	instructionType = "thawAccount"
	//	parsedInfo = map[string]interface{}{
	//		"account": ins.Accounts[0].PubKey.ToBase58(),
	//		"mint":    ins.Accounts[1].PubKey.ToBase58(),
	//	}
	//	parsedInfo = parse_signers(parsedInfo, 2, ins.Accounts, "freezeAuthority", "multisigFreezeAuthority")
	//
	//	break
	case InstructionTransferChecked:
		var a TransferCheckedInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		if err != nil {
			log.Printf("error unmarshalling InstructionTransferChecked: %v", err)
		}
		instructionType = "transferChecked"
		parsedInfo = map[string]interface{}{
			"source":      ins.Accounts[0].PubKey.ToBase58(),
			"mint":        ins.Accounts[1].PubKey.ToBase58(),
			"destination": ins.Accounts[2].PubKey.ToBase58(),
			"tokenAmount": tokenAmountToUiAmount(a.Amount, a.Decimals),
		}
		parsedInfo = parse_signers(parsedInfo, 3, ins.Accounts, "authority", "multisigAuthority")

		break
	//case InstructionApproveChecked:
	//	var a ApproveCheckedInstruction
	//	err = binstruct.UnmarshalLE(ins.Data, &a)
	//	instructionType = "approveChecked"
	//	parsedInfo = map[string]interface{}{
	//		"account":     ins.Accounts[0].PubKey.ToBase58(),
	//		"mint":        ins.Accounts[1].PubKey.ToBase58(),
	//		"delegate":    ins.Accounts[2].PubKey.ToBase58(),
	//		"tokenAmount": tokenAmountToUiAmount(a.Amount, a.Decimals),
	//	}
	//	parsedInfo = parse_signers(parsedInfo, 2, ins.Accounts, "owner", "owner")
	//
	//	break
	//case InstructionMintToChecked:
	//	var a MintToCheckedInstruction
	//	err = binstruct.UnmarshalLE(ins.Data, &a)
	//	instructionType = "mintToChecked"
	//	parsedInfo = map[string]interface{}{
	//		"account":     ins.Accounts[0].PubKey.ToBase58(),
	//		"mint":        ins.Accounts[1].PubKey.ToBase58(),
	//		"tokenAmount": tokenAmountToUiAmount(a.Amount, a.Decimals),
	//	}
	//	parsedInfo = parse_signers(parsedInfo, 2, ins.Accounts, "authority", "multisigAuthority")
	//
	//	break
	//case InstructionBurnChecked:
	//	var a BurnCheckedInstruction
	//	err = binstruct.UnmarshalLE(ins.Data, &a)
	//	instructionType = "burnChecked"
	//	parsedInfo = map[string]interface{}{
	//		"account":     ins.Accounts[0].PubKey.ToBase58(),
	//		"mint":        ins.Accounts[1].PubKey.ToBase58(),
	//		"tokenAmount": tokenAmountToUiAmount(a.Amount, a.Decimals),
	//	}
	//	parsedInfo = parse_signers(parsedInfo, 2, ins.Accounts, "authority", "multisigAuthority")
	//
	//	break
	default:
		log.Printf("ERROR unknown type='%v'", s.Instruction)
	}

	parsedInstruction.Parsed = &InstructionInfo{
		Info:            parsedInfo,
		InstructionType: instructionType,
	}
	return parsedInstruction, err
}

func parse_signers(
	m map[string]interface{},
	lastNonsignerIndex uint,
	accounts []types.AccountMeta,
	ownerFieldName string,
	multisigFieldName string,
) map[string]interface{} {
	if len(accounts) > int(lastNonsignerIndex)+1 {
		var signers []string
		for _, v := range accounts[lastNonsignerIndex+1:] {
			signers = append(signers, v.PubKey.ToBase58())
		}
		m[multisigFieldName] = accounts[lastNonsignerIndex].PubKey.ToBase58()
		m["signers"] = signers
	} else {
		m[ownerFieldName] = accounts[lastNonsignerIndex].PubKey.ToBase58()

	}
	return m
}

type UiTokenAmount struct {
	UiAmount float64
	Decimals uint8
	Amount   string
}

func tokenAmountToUiAmount(amount uint64, decimals uint8) UiTokenAmount {
	// Use `amount_to_ui_amount()` once spl_token is bumped to a version that supports it: https://github.com/solana-labs/solana-program-library/pull/211
	amountDecimals := float64(amount) / math.Pow(10, float64(decimals))
	uitokenAmt := UiTokenAmount{
		UiAmount: amountDecimals,
		Decimals: decimals,
		Amount:   fmt.Sprint(amount),
	}
	return uitokenAmt
}
