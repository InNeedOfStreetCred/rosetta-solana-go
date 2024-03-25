package operations

import (
	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/common"
	solPTypes "github.com/portto/solana-go-sdk/types"
	"log"
)

func AddSetComputeUnitPriceParam(microLamportsUnitPrice uint64, ins []solPTypes.Instruction) []solPTypes.Instruction {

	if microLamportsUnitPrice > 0 {
		log.Printf("Adding SetComputeUnitPriceParam=%v", microLamportsUnitPrice)
		ins = append(ins, SetComputeUnitPrice(SetComputeUnitPriceParam{MicroLamports: microLamportsUnitPrice}))
	} else {
		log.Printf("skipping SetComputeUnitPriceParam as fee is 0")
	}
	return ins
}

type SetComputeUnitLimitParam struct {
	Units uint32
}

type Instruction borsh.Enum

const (
	InstructionRequestUnits Instruction = iota
	InstructionRequestHeapFrame
	InstructionSetComputeUnitLimit
	InstructionSetComputeUnitPrice
)

func SetComputeUnitLimit(param SetComputeUnitLimitParam) solPTypes.Instruction {
	data, err := borsh.Serialize(struct {
		Instruction Instruction
		Units       uint32
	}{
		Instruction: InstructionSetComputeUnitLimit,
		Units:       param.Units,
	})
	if err != nil {
		panic(err)
	}

	return solPTypes.Instruction{
		ProgramID: common.PublicKeyFromString("ComputeBudget111111111111111111111111111111"),
		Accounts:  []solPTypes.AccountMeta{},
		Data:      data,
	}
}

type SetComputeUnitPriceParam struct {
	MicroLamports uint64
}

// SetComputeUnitPrice set a compute unit price in "micro-lamports" to pay a higher transaction
// fee for higher transaction prioritization.
func SetComputeUnitPrice(param SetComputeUnitPriceParam) solPTypes.Instruction {
	data, err := borsh.Serialize(struct {
		Instruction   Instruction
		MicroLamports uint64
	}{
		Instruction:   InstructionSetComputeUnitPrice,
		MicroLamports: param.MicroLamports,
	})
	if err != nil {
		panic(err)
	}

	return solPTypes.Instruction{
		ProgramID: common.PublicKeyFromString("ComputeBudget111111111111111111111111111111"),
		Accounts:  []solPTypes.AccountMeta{},
		Data:      data,
	}
}
