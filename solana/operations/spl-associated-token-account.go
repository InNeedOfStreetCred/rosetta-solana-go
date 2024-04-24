package operations

import (
	"encoding/json"
	"github.com/blocto/solana-go-sdk/common"
	token "github.com/blocto/solana-go-sdk/program/associated_token_account"
	solPTypes "github.com/blocto/solana-go-sdk/types"
	"github.com/coinbase/rosetta-sdk-go/types"
	stypes "github.com/imerkle/rosetta-solana-go/solana/shared_types"
)

type SplAssociatedTokenAccountOperationMetadata struct {
	Source string `json:"source,omitempty"`
	Wallet string `json:"wallet,omitempty"`
	Mint   string `json:"mint,omitempty"`
}

func (x *SplAssociatedTokenAccountOperationMetadata) SetMeta(op *types.Operation) {
	jsonString, _ := json.Marshal(op.Metadata)
	//if x.Source == "" {
	//	x.Source = op.Account.Address
	//}
	json.Unmarshal(jsonString, &x)
}

func (x *SplAssociatedTokenAccountOperationMetadata) ToInstructions(opType string) []solPTypes.Instruction {
	assosiatedAccount, _, _ := common.FindAssociatedTokenAddress(p(x.Wallet), p(x.Mint))
	var ins []solPTypes.Instruction
	switch opType {
	case stypes.SplAssociatedTokenAccount__Create:
		ins = append(ins, token.Create(token.CreateParam{Funder: p(x.Source), Owner: p(x.Wallet), Mint: p(x.Mint), AssociatedTokenAccount: assosiatedAccount}))
		break
	}
	return ins
}
