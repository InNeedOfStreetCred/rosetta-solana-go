package operations

import (
	"encoding/json"
	types "github.com/coinbase/rosetta-sdk-go/types"
	solanago "github.com/imerkle/rosetta-solana-go/solana"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/stakeprog"
	"github.com/portto/solana-go-sdk/sysprog"
	solPTypes "github.com/portto/solana-go-sdk/types"
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

func (x *StakeOperationMetadata) SetMeta(op *types.Operation, fee solanago.PriorityFee) {
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
	case solanago.Stake__CreateStakeAccount:
		ins = addCreateStakeAccountIns(ins, x)
		break
	case solanago.Stake__DelegateStake:
		ins = addDelegateStakeIns(ins, x)
		break
	case solanago.Stake__CreateStakeAndDelegate:
		ins = addCreateStakeAccountIns(ins, x)
		ins = addDelegateStakeIns(ins, x)
		break
	case solanago.Stake__DeactivateStake:
		ins = append(ins, stakeprog.Deactivate(p(x.Stake), p(x.Staker)))
		break
	case solanago.Stake__WithdrawStake:
		ins = append(ins,
			stakeprog.Withdraw(
				p(x.Stake),
				p(x.Withdrawer),
				p(x.WithdrawDestination),
				x.Lamports,
				p(x.LockupCustodian)))
		break
	case solanago.Stake__Merge:
		ins = append(ins,
			stakeprog.Merge(
				p(x.MergeDestination),
				p(x.Stake),
				p(x.Staker)))
		break
	case solanago.Stake__Split:
		ins = append(ins,
			stakeprog.Split(
				p(x.Stake),
				p(x.Staker),
				p(x.SplitDestination),
				x.Lamports))
		break
	case solanago.Stake__Authorize:
		ins = append(ins,
			stakeprog.Authorize(
				p(x.Stake),
				p(x.Authority),
				p(x.NewAuthority),
				stakeprog.StakeAuthorizationType(x.StakeAuthorizationType),
				p(x.LockupCustodian)))
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

	log.Printf("END ToInstructions")
	return ins
}

func addCreateStakeAccountIns(ins []solPTypes.Instruction, x *StakeOperationMetadata) []solPTypes.Instruction {
	ins = append(ins,
		sysprog.CreateAccount(
			p(x.Source),
			p(x.Stake),
			common.StakeProgramID,
			x.Lamports,
			stakeprog.AccountSize))
	ins = append(ins,
		stakeprog.Initialize(
			p(x.Stake),
			stakeprog.Authorized{
				Staker:     p(x.Staker),
				Withdrawer: p(x.Withdrawer),
			},
			stakeprog.Lockup{
				UnixTimestamp: x.LockupUnixTimestamp,
				Epoch:         x.LockupEpoch,
				Cusodian:      p(x.LockupCustodian),
			}))
	return ins
}

func addDelegateStakeIns(ins []solPTypes.Instruction, x *StakeOperationMetadata) []solPTypes.Instruction {
	return append(ins, stakeprog.DelegateStake(p(x.Stake), p(x.Staker), p(x.VoteAccount)))
}
