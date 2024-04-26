package parse

import (
	"github.com/blocto/solana-go-sdk/common"
	types "github.com/blocto/solana-go-sdk/types"
	"github.com/imerkle/rosetta-solana-go/solana/parse/associatedtokenaccount"
	"github.com/imerkle/rosetta-solana-go/solana/parse/computebudget"
	"github.com/imerkle/rosetta-solana-go/solana/parse/system"
	"github.com/imerkle/rosetta-solana-go/solana/parse/token"
	stypes "github.com/imerkle/rosetta-solana-go/solana/shared_types"
	"github.com/mr-tron/base58"
	"log"
)

var (
	SystemProgramID                    = common.PublicKeyFromString("11111111111111111111111111111111")
	ConfigProgramID                    = common.PublicKeyFromString("Config1111111111111111111111111111111111111")
	StakeProgramID                     = common.PublicKeyFromString("Stake11111111111111111111111111111111111111")
	VoteProgramID                      = common.PublicKeyFromString("Vote111111111111111111111111111111111111111")
	BPFLoaderProgramID                 = common.PublicKeyFromString("BPFLoader1111111111111111111111111111111111")
	Secp256k1ProgramID                 = common.PublicKeyFromString("KeccakSecp256k11111111111111111111111111111")
	TokenProgramID                     = common.PublicKeyFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	SPLAssociatedTokenAccountProgramID = common.PublicKeyFromString("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")
	ComputeBudgetProgramID             = common.PublicKeyFromString("ComputeBudget111111111111111111111111111111")
)

func ToParsedTransaction(tx types.Transaction) (stypes.ParsedTransaction, error) {
	ins := tx.Message.DecompileInstructions()
	var parsedIns []stypes.ParsedInstruction
	for _, v := range ins {
		p, err := ParseInstruction(v)
		if err != nil {
			//cannot parse
			p = stypes.ParsedInstruction{}
			return stypes.ParsedTransaction{}, err
		}
		parsedIns = append(parsedIns, p)
	}
	var acckeys []stypes.ParsedAccKey
	var sigs []string
	for _, v := range tx.Message.Accounts {
		acckeys = append(acckeys, stypes.ParsedAccKey{PubKey: v.ToBase58()})
	}
	for _, v := range tx.Signatures {
		sigs = append(sigs, base58.Encode(v[:]))
	}
	newTx := stypes.ParsedTransaction{
		Signatures: sigs,
		Message: stypes.ParsedMessage{
			Header:          tx.Message.Header,
			AccountKeys:     acckeys,
			RecentBlockhash: tx.Message.RecentBlockHash,
			Instructions:    parsedIns,
		},
	}
	return newTx, nil
}

func ParseInstruction(ins types.Instruction) (stypes.ParsedInstruction, error) {
	var parsedInstruction stypes.ParsedInstruction
	var err error

	switch ins.ProgramID {
	case common.SystemProgramID:
		parsedInstruction, err = system.ParseSystem(ins)
		if err != nil {
			log.Printf("error parsing SystemProgramID instruction: %v", err)
		}
		break
	case common.TokenProgramID:
		parsedInstruction, err = token.ParseToken(ins)
		if err != nil {
			log.Printf("error parsing TokenProgramID instruction: %v", err)
		}
		break
	case common.SPLAssociatedTokenAccountProgramID:
		parsedInstruction, err = associatedtokenaccount.ParseAssocToken(ins)
		if err != nil {
			log.Printf("error parsing SPLAssociatedTokenAccountProgramID instruction: %v", err)
		}
		break

	case common.ComputeBudgetProgramID:
		parsedInstruction, err = computebudget.ParseComputeBudget(ins)
		if err != nil {
			log.Printf("error parsing ComputeBudgetProgramID instruction: %v", err)
		}
		break
	default:
		log.Printf("error parsing instruction. ins.ProgramI=%v is unknown", ins.ProgramID)
		//return parsedInstruction, fmt.Errorf("Cannot parse instruction")
	}
	if err != nil {
		log.Printf("error parsing instruction: %v", err)
		return parsedInstruction, err
	}
	var accs []string
	for _, v := range ins.Accounts {
		accs = append(accs, v.PubKey.ToBase58())
	}
	parsedInstruction.Accounts = accs
	parsedInstruction.Data = base58.Encode(ins.Data[:])
	parsedInstruction.ProgramID = ins.ProgramID.ToBase58()
	parsedInstruction.Program = GetProgramName(ins.ProgramID)
	return parsedInstruction, nil
}

func GetProgramName(programId common.PublicKey) string {
	name := "Unknown"
	switch programId {
	case SystemProgramID:
		name = "system"
		break
	case ConfigProgramID:
		name = "spl-token"
		break
	case StakeProgramID:
		name = "stake"
		break
	case VoteProgramID:
		name = "vote"
		break
	case BPFLoaderProgramID:
		name = "bpf-loader"
		break
	case Secp256k1ProgramID:
		name = "secp256k1"
		break
	case TokenProgramID:
		name = "spl-token"
		break
	case SPLAssociatedTokenAccountProgramID:
		name = "spl-associated-token-account"
		break
	case ComputeBudgetProgramID:
		name = "compute-budget"
		break
	}
	return name
}
