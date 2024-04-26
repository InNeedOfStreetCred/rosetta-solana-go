package computebudget

import (
	"github.com/blocto/solana-go-sdk/types"
	"github.com/ghostiam/binstruct"
	stypes "github.com/imerkle/rosetta-solana-go/solana/shared_types"
	"log"
)

type Instruction uint8

const (
	InstructionRequestUnits Instruction = iota
	InstructionRequestHeapFrame
	InstructionSetComputeUnitLimit
	InstructionSetComputeUnitPrice
)

func ParseComputeBudget(ins types.Instruction) (stypes.ParsedInstruction, error) {

	var parsedInstruction stypes.ParsedInstruction
	var err error
	var s struct {
		Instruction Instruction
	}
	err = binstruct.UnmarshalLE(ins.Data, &s)
	var instructionType string
	var parsedInfo map[string]interface{}
	switch s.Instruction {
	case InstructionSetComputeUnitPrice:
		var a SetComputeUnitPriceInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "setComputeUnitPrice"
		parsedInfo = map[string]interface{}{
			"MicroLamports": a.MicroLamports,
		}
		break
	}

	if err != nil {
		log.Printf("error parsing instruction: %v", err)
		return parsedInstruction, err
	}

	parsedInstruction.Parsed = &stypes.InstructionInfo{
		Info:            parsedInfo,
		InstructionType: instructionType,
	}
	return parsedInstruction, nil
}

type SetComputeUnitPriceInstruction struct {
	Instruction   Instruction
	MicroLamports uint64
}
