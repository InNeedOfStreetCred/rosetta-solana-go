package operations

import (
	"encoding/json"
	"github.com/blocto/solana-go-sdk/program/stake"
	"github.com/blocto/solana-go-sdk/program/stakeprog"
	"github.com/blocto/solana-go-sdk/program/system"
	solPTypes "github.com/blocto/solana-go-sdk/types"
	types "github.com/coinbase/rosetta-sdk-go/types"
	solanago "github.com/imerkle/rosetta-solana-go/solana"
	stypes "github.com/imerkle/rosetta-solana-go/solana/shared_types"
	"log"
)

type StakeOperationMetadata struct {
	Source                 string `json:"source,omitempty"`
	Stake                  string `json:"stake,omitempty"`
	Lamports               uint64 `json:"lamports,omitempty"`
	Staker                 string `json:"staker,omitempty"`
	Withdrawer             string `json:"withdrawer,omitempty"`
	WithdrawDestination    string `json:"withdrawDestination,omitempty"`
	LockupUnixTimestamp    int64  `json:"lockupUnixTimestamp,omitempty"`
	LockupEpoch            uint64 `json:"lockupEpoch,omitempty"`
	LockupCustodian        string `json:"lockupCustodian,omitempty"`
	VoteAccount            string `json:"voteAccount,omitempty"`
	MergeDestination       string `json:"mergeDestination,omitempty"`
	SplitDestination       string `json:"splitDestination,omitempty"`
	Authority              string `json:"authority,omitempty"`
	NewAuthority           string `json:"newAuthority,omitempty"`
	StakeAuthorizationType uint32 `json:"stakeAuthorizationType,omitempty"`
	FeePayer               string `json:"feePayer,omitempty"`
	MicroLamportsUnitPrice uint64 `json:"micro_lamports_unit_price,omitempty"`
}

func (x *StakeOperationMetadata) SetMeta(op *types.Operation, fee stypes.PriorityFee) {
	jsonString, _ := json.Marshal(op.Metadata)
	json.Unmarshal(jsonString, &x)
	if x.Lamports == 0 && op.Amount != nil {
		x.Lamports = solanago.ValueToBaseAmount(op.Amount.Value)
	}
	if x.Source == "" && op.Account != nil {
		x.Source = op.Account.Address
	}
	if x.Staker == "" && op.Account != nil {
		x.Staker = op.Account.Address
	}
	if x.Withdrawer == "" && op.Account != nil {
		x.Withdrawer = op.Account.Address
	}
	if fee.MicroLamports != "" {
		x.MicroLamportsUnitPrice = solanago.ValueToBaseAmount(fee.MicroLamports)
	}
	log.Printf("microLamportsUnitPrice=%v", x.MicroLamportsUnitPrice)
}
func (x *StakeOperationMetadata) ToInstructions(opType string) []solPTypes.Instruction {
	log.Printf("START ToInstructions")
	log.Printf("opType=%v", opType)

	var ins []solPTypes.Instruction
	ins = AddSetComputeUnitPriceParam(x.MicroLamportsUnitPrice, ins)
	switch opType {
	case stypes.Stake__CreateStakeAccount:
		ins = addCreateStakeAccountIns(ins, x)
		break
	case stypes.Stake__DelegateStake:
		ins = addDelegateStakeIns(ins, x)
		break
	case stypes.Stake__CreateStakeAndDelegate:
		ins = addCreateStakeAccountIns(ins, x)
		ins = addDelegateStakeIns(ins, x)
		break
	case stypes.Stake__DeactivateStake:
		ins = append(ins, stake.Deactivate(stake.DeactivateParam{Stake: p(x.Stake), Auth: p(x.Staker)}))
		break
	case stypes.Stake__WithdrawStake:
		lockupCustodian := p(x.LockupCustodian)
		ins = append(ins,
			stake.Withdraw(
				stake.WithdrawParam{
					Stake:     p(x.Stake),
					Auth:      p(x.Withdrawer),
					To:        p(x.WithdrawDestination),
					Lamports:  x.Lamports,
					Custodian: &lockupCustodian}))
		break
		//case solanago.Stake__Merge:
		//	ins = append(ins,
		//		stake.Merge(
		//			p(x.MergeDestination),
		//			p(x.Stake),
		//			p(x.Staker)))
		//	break
		//case solanago.Stake__Split:
		//	ins = append(ins,
		//		stake.Split(
		//			p(x.Stake),
		//			p(x.Staker),
		//			p(x.SplitDestination),
		//			x.Lamports))
		//	break
		//case solanago.Stake__Authorize:
		//	ins = append(ins,
		//		stake.Authorize(
		//			p(x.Stake),
		//			p(x.Authority),
		//			p(x.NewAuthority),
		//			stakeprog.StakeAuthorizationType(x.StakeAuthorizationType),
		//			p(x.LockupCustodian)))
		//	break
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

	log.Printf("END ToInstructions")
	return ins
}

func addCreateStakeAccountIns(ins []solPTypes.Instruction, x *StakeOperationMetadata) []solPTypes.Instruction {
	ins = append(ins,
		system.CreateAccount(
			system.CreateAccountParam{
				From:     p(x.Source),
				New:      p(x.Stake),
				Lamports: x.Lamports,
				Space:    stakeprog.AccountSize}))
	ins = append(ins,
		stake.Initialize(
			stake.InitializeParam{
				Stake: p(x.Stake),
				Auth: stake.Authorized{
					Staker:     p(x.Staker),
					Withdrawer: p(x.Withdrawer),
				},
				Lockup: stake.Lockup{
					UnixTimestamp: x.LockupUnixTimestamp,
					Epoch:         x.LockupEpoch,
					Cusodian:      p(x.LockupCustodian),
				}}),
	)
	return ins
}

func addDelegateStakeIns(ins []solPTypes.Instruction, x *StakeOperationMetadata) []solPTypes.Instruction {
	return append(ins, stake.DelegateStake(stake.DelegateStakeParam{Stake: p(x.Stake), Auth: p(x.Staker), Vote: p(x.VoteAccount)}))
}
