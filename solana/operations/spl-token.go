package operations

import (
	"encoding/json"
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/program/tokenprog"
	solPTypes "github.com/blocto/solana-go-sdk/types"
	"github.com/coinbase/rosetta-sdk-go/types"
	solanago "github.com/imerkle/rosetta-solana-go/solana"
)

type SplTokenOperationMetadata struct {
	Source          string `json:"source,omitempty"`
	Destination     string `json:"destination,omitempty"`
	Mint            string `json:"mint,omitempty"`
	Authority       string `json:"authority,omitempty"`
	FreezeAuthority string `json:"freeze_authority,omitempty"`
	Amount          uint64 `json:"amount,omitempty"`
	Decimals        uint8  `json:"decimals,omitempty"`

	SourceToken      string `json:"source_token,omitempty"`
	DestinationToken string `json:"destination_token,omitempty"`
}

func (x *SplTokenOperationMetadata) SetMeta(op *types.Operation, splTokenAccsMap map[string]solanago.SplAccounts) {
	jsonString, _ := json.Marshal(op.Metadata)
	if op.Amount != nil && x.Amount == 0 {
		x.Amount = solanago.ValueToBaseAmount(op.Amount.Value)
	}
	if x.Source == "" {
		x.Source = op.Account.Address
	}
	if x.Authority == "" {
		x.Authority = x.Source
	}
	if op.Amount != nil && x.Mint == "" {
		x.Mint = op.Amount.Currency.Symbol
	}
	if op.Amount != nil && x.Decimals == 0 {
		x.Decimals = uint8(op.Amount.Currency.Decimals)
	}
	if w, ok := splTokenAccsMap[fmt.Sprint(op.OperationIdentifier.Index)]; ok {
		x.SourceToken = w.Source
		x.DestinationToken = w.Destination
	}

	json.Unmarshal(jsonString, &x)
}

func (x *SplTokenOperationMetadata) ToInstructions(opType string) []solPTypes.Instruction {

	var ins []solPTypes.Instruction
	switch opType {
	//	case solanago.SplToken__InitializeMint:
	//		ins = append(ins, tokenprog.InitializeMint(x.Decimals, p(x.Mint), p(x.Source), p(x.Authority)))
	//		break
	//	case solanago.SplToken__InitializeAccount:
	//		ins = append(ins, tokenprog.InitializeAccount(p(x.Destination), p(x.Mint), p(x.Source)))
	//
	//		break
	//	case solanago.SplToken__CreateToken:
	//		ins = append(ins, sysprog.CreateAccount(p(x.Source), p(x.Mint), common.TokenProgramID, x.Amount, tokenprog.MintAccountSize))
	//		ins = append(ins, tokenprog.InitializeMint(x.Decimals, p(x.Mint), p(x.Source), p(x.Authority)))
	//		break
	case solanago.SplToken__CreateAccount:
		ins = append(ins, system.CreateAccount(system.CreateAccountParam{From: p(x.Source), New: p(x.Destination), Lamports: x.Amount, Space: tokenprog.TokenAccountSize}))
		ins = append(ins, token.InitializeAccount(token.InitializeAccountParam{Account: p(x.Destination), Mint: p(x.Mint), Owner: p(x.Authority)}))

		break
		//	case solanago.SplToken__Approve:
		//		ins = append(ins, tokenprog.Approve(p(x.Source), p(x.Destination), p(x.Authority), []common.PublicKey{}, x.Amount))
		//		break
		//	case solanago.SplToken__Revoke:
		//		ins = append(ins, tokenprog.Revoke(p(x.Source), p(x.Authority), []common.PublicKey{}))
		//		break
		//	case solanago.SplToken_MintTo:
		//		ins = append(ins, tokenprog.MintTo(p(x.Mint), p(x.Source), p(x.Authority), []common.PublicKey{}, x.Amount))
		//		break
		//	case solanago.SplToken_Burn:
		//		ins = append(ins, tokenprog.Burn(p(x.Source), p(x.Mint), p(x.Authority), []common.PublicKey{}, x.Amount))
		//		break
		//	case solanago.SplToken_CloseAccount:
		//		ins = append(ins, tokenprog.CloseAccount(p(x.Source), p(x.Destination), p(x.Authority), []common.PublicKey{}))
		//		break
		//	case solanago.SplToken_FreezeAccount:
		//		ins = append(ins, tokenprog.ThawAccount(p(x.Source), p(x.Mint), p(x.Authority), []common.PublicKey{}))
		//		break
		//	case solanago.SplToken__Transfer:
		//		ins = append(ins, tokenprog.Transfer(p(x.Source), p(x.Destination), p(x.Authority), []common.PublicKey{}, x.Amount))
		//		break
		//	case solanago.SplToken__TransferChecked:
		//		ins = append(ins, tokenprog.TransferChecked(p(x.Source), p(x.Destination), p(x.Mint), p(x.Authority), []common.PublicKey{}, x.Amount, x.Decimals))
		//		break
		//	case solanago.SplToken__TransferNew:
		//		ins_create_assoc := assotokenprog.CreateAssociatedTokenAccount(p(x.Authority), p(x.Destination), p(x.Mint))
		//		account := ins_create_assoc.Accounts[1].PubKey.ToBase58()
		//		ins = append(ins, ins_create_assoc)
		//		ins = append(ins, tokenprog.TransferChecked(p(x.Source), p(account), p(x.Mint), p(x.Authority), []common.PublicKey{}, x.Amount, x.Decimals))
		//		break
		//	case solanago.SplToken__TransferWithSystem:
		//		source := x.SourceToken
		//		destination := x.DestinationToken
		//		if x.SourceToken == "" {
		//			in := assotokenprog.CreateAssociatedTokenAccount(p(x.Authority), p(x.Source), p(x.Mint))
		//			source = in.Accounts[1].PubKey.ToBase58()
		//			ins = append(ins, in)
		//		}
		//		if x.DestinationToken == "" {
		//			in := assotokenprog.CreateAssociatedTokenAccount(p(x.Authority), p(x.Destination), p(x.Mint))
		//			destination = in.Accounts[1].PubKey.ToBase58()
		//			ins = append(ins, in)
		//		}
		//		ins = append(ins, tokenprog.TransferChecked(p(source), p(destination), p(x.Mint), p(x.Authority), []common.PublicKey{}, x.Amount, x.Decimals))
		//		break
	}
	return ins
}

func p(a string) common.PublicKey {
	return common.PublicKeyFromString(a)
}
