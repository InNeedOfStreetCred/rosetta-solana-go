package operations

import (
	"encoding/json"
	"fmt"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/associated_token_account"
	"github.com/blocto/solana-go-sdk/program/assotokenprog"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/blocto/solana-go-sdk/program/token"
	"github.com/blocto/solana-go-sdk/program/tokenprog"
	solPTypes "github.com/blocto/solana-go-sdk/types"
	"github.com/coinbase/rosetta-sdk-go/types"
	solanago "github.com/imerkle/rosetta-solana-go/solana"
	stypes "github.com/imerkle/rosetta-solana-go/solana/shared_types"
	"log"
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

func (x *SplTokenOperationMetadata) SetMeta(op *types.Operation, splTokenAccsMap map[string]stypes.SplAccounts) {
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
	case stypes.SplToken__CreateAccount:
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
	case stypes.SplToken__Transfer:
		param := token.TransferParam{
			From:    p(x.Source),
			To:      p(x.Destination),
			Auth:    p(x.Authority),
			Signers: []common.PublicKey{},
			Amount:  x.Amount}
		ins = append(ins, token.Transfer(param))
		break
	case stypes.SplToken__TransferChecked:
		param := token.TransferCheckedParam{
			From:     p(x.Source),
			To:       p(x.Destination),
			Mint:     p(x.Mint),
			Auth:     p(x.Authority),
			Signers:  []common.PublicKey{},
			Amount:   x.Amount,
			Decimals: x.Decimals}
		ins = append(ins, token.TransferChecked(param))
		break
	case stypes.SplToken__TransferNew:
		assosiatedAccount, _, _ := common.FindAssociatedTokenAddress(p(x.Destination), p(x.Mint))
		ins_create_assoc := assotokenprog.CreateAssociatedTokenAccount(assotokenprog.CreateAssociatedTokenAccountParam{Funder: p(x.Authority), Owner: p(x.Destination), Mint: p(x.Mint), AssociatedTokenAccount: assosiatedAccount})
		account := ins_create_assoc.Accounts[1].PubKey.ToBase58()
		ins = append(ins, ins_create_assoc)
		ins = append(ins, tokenprog.TransferChecked(tokenprog.TransferCheckedParam{From: p(x.Source), To: p(account), Mint: p(x.Mint), Auth: p(x.Authority), Signers: []common.PublicKey{}, Amount: x.Amount, Decimals: x.Decimals}))
		break
	case stypes.SplToken__TransferWithSystem:
		source := x.SourceToken
		destination := x.DestinationToken
		if x.SourceToken == "" {
			assosiatedAccount, _, _ := common.FindAssociatedTokenAddress(p(source), p(x.Mint))
			param := associated_token_account.CreateIdempotentParam{Funder: p(x.Authority), Owner: p(source), Mint: p(x.Mint), AssociatedTokenAccount: assosiatedAccount}
			in := associated_token_account.CreateIdempotent(param)
			source = in.Accounts[1].PubKey.ToBase58()
			ins = append(ins, in)
		}
		if x.DestinationToken == "" {
			assosiatedAccount, _, _ := common.FindAssociatedTokenAddress(p(x.Destination), p(x.Mint))
			param := associated_token_account.CreateIdempotentParam{Funder: p(x.Authority), Owner: p(x.Destination), Mint: p(x.Mint), AssociatedTokenAccount: assosiatedAccount}
			in := associated_token_account.CreateIdempotent(param)
			destination = in.Accounts[1].PubKey.ToBase58()
			ins = append(ins, in)
		}
		ins = append(ins, tokenprog.TransferChecked(tokenprog.TransferCheckedParam{From: p(x.Source), To: p(destination), Mint: p(x.Mint), Auth: p(x.Authority), Signers: []common.PublicKey{}, Amount: x.Amount, Decimals: x.Decimals}))
		break
	default:
		log.Printf("ERROR: unknown opType='%v'", opType)
	}
	return ins
}

func p(a string) common.PublicKey {
	return common.PublicKeyFromString(a)
}
