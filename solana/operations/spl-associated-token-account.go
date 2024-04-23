package operations

import (
	"encoding/json"

	"github.com/coinbase/rosetta-sdk-go/types"
)

type SplAssociatedTokenAccountOperationMetadata struct {
	Source string `json:"source,omitempty"`
	Wallet string `json:"wallet,omitempty"`
	Mint   string `json:"mint,omitempty"`
}

func (x *SplAssociatedTokenAccountOperationMetadata) SetMeta(op *types.Operation) {
	jsonString, _ := json.Marshal(op.Metadata)
	if x.Source == "" {
		x.Source = op.Account.Address
	}
	json.Unmarshal(jsonString, &x)
}

//func (x *SplAssociatedTokenAccountOperationMetadata) ToInstructions(opType string) []solPTypes.Instruction {
//	var ins []solPTypes.Instruction
//	switch opType {
//	case solanago.SplAssociatedTokenAccount__Create:
//		ins = append(ins, assotokenprog.Create(p(x.Source), p(x.Wallet), p(x.Mint)))
//		break
//	}
//	return ins
//}
