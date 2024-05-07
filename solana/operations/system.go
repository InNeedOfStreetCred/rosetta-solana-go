package operations

import (
	"encoding/json"
	"github.com/blocto/solana-go-sdk/common"
	"github.com/blocto/solana-go-sdk/program/sysprog"
	"github.com/blocto/solana-go-sdk/program/system"
	solPTypes "github.com/blocto/solana-go-sdk/types"
	"github.com/coinbase/rosetta-sdk-go/types"
	solanago "github.com/imerkle/rosetta-solana-go/solana"
	stypes "github.com/imerkle/rosetta-solana-go/solana/shared_types"
	"log"
)

type SystemOperationMetadata struct {
	Source                 string `json:"source,omitempty"`
	Destination            string `json:"destination,omitempty"`
	Space                  uint64 `json:"space,omitempty"`
	Lamports               uint64 `json:"lamports,omitempty"`
	NewAuthority           string `json:"new_authority,omitempty"`
	Authority              string `json:"authority,omitempty"`
	MicroLamportsUnitPrice uint64 `json:"micro_lamports_unit_price,omitempty"`
}

func (x *SystemOperationMetadata) SetMeta(op *types.Operation, fee stypes.PriorityFee) {
	jsonString, _ := json.Marshal(op.Metadata)
	json.Unmarshal(jsonString, &x)
	if x.Lamports == 0 {
		x.Lamports = solanago.ValueToBaseAmount(op.Amount.Value)
	}
	if x.Source == "" {
		x.Source = op.Account.Address
	}
	if x.Authority == "" {
		x.Authority = x.Source
	}
	if fee.MicroLamports != "" {
		x.MicroLamportsUnitPrice = solanago.ValueToBaseAmount(fee.MicroLamports)
	}
	log.Printf("microLamportsUnitPrice=%v", x.MicroLamportsUnitPrice)
}
func (x *SystemOperationMetadata) ToInstructions(opType string) []solPTypes.Instruction {
	log.Printf("START system ToInstructions")
	log.Printf("opType=%v", opType)

	var ins []solPTypes.Instruction
	ins = AddSetComputeUnitPriceParam(x.MicroLamportsUnitPrice, ins)
	switch opType {
	case stypes.System__CreateAccount:
		log.Printf("System__CreateAccount adding CreateAccount")
		ins = append(ins, system.CreateAccount(system.CreateAccountParam{From: p(x.Source), New: p(x.Destination), Lamports: x.Lamports, Space: x.Space}))
		break
	case stypes.System__Assign:
		ins = append(ins, system.Assign(system.AssignParam{From: p(x.Source)}))
		break
	case stypes.System__Transfer:

		log.Printf("System__Transfer adding transfer")
		ins = append(ins, system.Transfer(system.TransferParam{From: p(x.Source), To: p(x.Destination), Amount: x.Lamports}))
		break
	case stypes.System__CreateNonceAccount:
		log.Printf("System__CreateNonceAccount adding CreateAccount")
		ins = append(ins, system.CreateAccount(system.CreateAccountParam{From: p(x.Source), New: p(x.Destination), Lamports: x.Lamports, Space: sysprog.NonceAccountSize}))
		log.Printf("System__CreateNonceAccount adding InitializeNonceAccount")
		ins = append(ins, solPTypes.Instruction{
			Accounts: []solPTypes.AccountMeta{
				{PubKey: p(x.Destination), IsSigner: false, IsWritable: true},
				{PubKey: common.SysVarRecentBlockhashsPubkey, IsSigner: false, IsWritable: false},
				{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
			},
			ProgramID: common.SystemProgramID,
			Data:      system.InitializeNonceAccount(system.InitializeNonceAccountParam{Nonce: p(x.Destination), Auth: p(x.Authority)}).Data,
		})

		break
	case stypes.System__AdvanceNonce:
		log.Printf("System__AdvanceNonce adding AdvanceNonceAccount")
		ins = append(ins, system.AdvanceNonceAccount(system.AdvanceNonceAccountParam{Nonce: p(x.Destination), Auth: p(x.Authority)}))
		break
	case stypes.System__WithdrawFromNonce:
		ins = append(ins, system.WithdrawNonceAccount(system.WithdrawNonceAccountParam{p(x.Source), p(x.Authority), p(x.Destination), x.Lamports}))
		break
	case stypes.System__AuthorizeNonce:
		ins = append(ins, system.AuthorizeNonceAccount(system.AuthorizeNonceAccountParam{p(x.Destination), p(x.Authority), p(x.NewAuthority)}))
		break
	case stypes.System__Allocate:
		ins = append(ins, system.Allocate(system.AllocateParam{p(x.Source), x.Space}))
		break
	}
	log.Printf("There are %v instructions", len(ins))
	for i, in := range ins {
		log.Printf("instruction with i=%v", i)
		log.Printf("in.ProgramID=%v", in.ProgramID.ToBase58())
		if (in.Accounts != nil) && (len(in.Accounts) > 0) {
			for _, account := range in.Accounts {
				log.Printf("account.PubKey=%v, IsSigner=%v, IsWritable=%v", account.PubKey.ToBase58(), account.IsSigner, account.IsWritable)
			}
		}
		log.Printf("in.Data=%v", in.Data)
	}

	log.Printf("END system ToInstructions")
	return ins
}
