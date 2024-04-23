package solanago

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	common "github.com/blocto/solana-go-sdk/common"
	solPTypes "github.com/blocto/solana-go-sdk/types"
	"github.com/coinbase/rosetta-sdk-go/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/ghostiam/binstruct"
	"github.com/mr-tron/base58"

	"github.com/iancoleman/strcase"
)

type InstructionInt uint32

const (
	InstructionCreateAccount InstructionInt = iota
	InstructionAssign
	InstructionTransfer
	InstructionCreateAccountWithSeed
	InstructionAdvanceNonceAccount
	InstructionWithdrawNonceAccount
	InstructionInitializeNonceAccount
	InstructionAuthorizeNonceAccount
	InstructionAllocate
	InstructionAllocateWithSeed
	InstructionAssignWithSeed
	InstructionTransferWithSeed
	InstructionUpgradeNonceAccount
)

func IsBalanceChanging(opType string) bool {
	a := false
	switch opType {
	case System__CreateAccount, System__WithdrawFromNonce, System__Transfer, SplToken__Transfer, SplToken__TransferChecked, "Stake__Split", "Stake__Withdraw", "Vote__Withdraw", SplToken__TransferNew, SplToken__TransferWithSystem:
		a = true
	}
	return a
}

func getOperationTypeWithProgram(program string, s string) string {
	toPascal := strcase.ToCamel(program)

	newStr := fmt.Sprint(
		toPascal,
		Separator,
		strcase.ToCamel(s),
	)
	return newStr
}
func getOperationType(s string) string {
	x := strings.Split(s, Separator)
	if len(x) < 2 {
		return Unknown
	}
	return getOperationTypeWithProgram(x[0], x[1])
}
func split_at(at int, input []byte) ([]byte, []byte) {
	return input[0:1], input[1:]
}
func GetRosOperationsFromTx(tx ParsedTransaction, status string) []*types.Operation {
	//	hash := tx.Transaction.Signatures[0].String()
	opIndex := int64(0)
	var operations []*types.Operation
	for _, ins := range tx.Message.Instructions {
		oi := types.OperationIdentifier{
			Index: opIndex,
		}
		opIndex += 1

		if ins.Parsed == nil {

			var inInterface map[string]interface{}
			inrec, _ := json.Marshal(ins)
			json.Unmarshal(inrec, &inInterface)

			operations = append(operations, &types.Operation{
				OperationIdentifier: &oi,
				Type:                Unknown,
				Status:              &status,
				Metadata:            inInterface,
			})
		} else {

			jsonString, _ := json.Marshal(ins.Parsed.Info)

			parsedInstructionMeta := ParsedInstructionMeta{}
			var parsedInstructionMetaInterface interface{}
			json.Unmarshal(jsonString, &parsedInstructionMeta)
			json.Unmarshal(jsonString, &parsedInstructionMetaInterface)

			var inInterface map[string]interface{}
			inrec, _ := json.Marshal(parsedInstructionMetaInterface)
			json.Unmarshal(inrec, &inInterface)

			opType := getOperationTypeWithProgram(ins.Program, ins.Parsed.InstructionType)
			if !Contains(OperationTypes, opType) {
				inInterface["instruction_type"] = ins.Parsed.InstructionType
				inInterface["program"] = ins.Program
				opType = "Unknown"
			}
			if IsBalanceChanging(opType) {
				if parsedInstructionMeta.Decimals == 0 {
					parsedInstructionMeta.Decimals = Decimals
				}
				if parsedInstructionMeta.Amount == 0 {
					if parsedInstructionMeta.Lamports == 0 {
						parsedInstructionMeta.Amount, _ = strconv.ParseUint(parsedInstructionMeta.TokenAmount.Amount, 10, 64)
					} else {
						parsedInstructionMeta.Amount = parsedInstructionMeta.Lamports
					}
				}
				var currency types.Currency
				if parsedInstructionMeta.Mint == "" {
					if ins.Program == "system" {
						currency = types.Currency{
							Symbol:   Symbol,
							Decimals: Decimals,
							Metadata: map[string]interface{}{},
						}
					}
				} else {
					currency = types.Currency{
						Symbol:   parsedInstructionMeta.Mint,
						Decimals: int32(parsedInstructionMeta.Decimals),
						Metadata: map[string]interface{}{},
					}
				}

				source := parsedInstructionMeta.Source
				if source == "" {
					source = parsedInstructionMeta.Owner
				}
				sender := types.AccountIdentifier{
					Address:  source,
					Metadata: map[string]interface{}{},
				}
				senderAmt := types.Amount{
					Value:    "-" + fmt.Sprint(parsedInstructionMeta.Amount),
					Currency: &currency,
				}

				destination := parsedInstructionMeta.Destination
				if destination == "" {
					destination = parsedInstructionMeta.NewAccount
				}
				receiver := types.AccountIdentifier{
					Address:  destination,
					Metadata: map[string]interface{}{},
				}
				receiverAmt := types.Amount{
					Value:    fmt.Sprint(parsedInstructionMeta.Amount),
					Currency: &currency,
				}
				oi2 := types.OperationIdentifier{
					Index: opIndex,
				}
				opIndex += 1

				//for construction test
				delete(inInterface, "amount")
				delete(inInterface, "lamports")
				delete(inInterface, "source")
				delete(inInterface, "destination")

				//sender push
				operations = append(operations, &types.Operation{
					OperationIdentifier: &oi,
					Type:                opType,
					Status:              &status,
					Account:             &sender,
					Amount:              &senderAmt,
					Metadata:            inInterface,
				}, &types.Operation{
					OperationIdentifier: &oi2,
					Type:                opType,
					Status:              &status,
					Account:             &receiver,
					Amount:              &receiverAmt,
					Metadata:            inInterface,
				})
			} else {
				var account types.AccountIdentifier
				if parsedInstructionMeta.Source != "" {
					account = types.AccountIdentifier{
						Address: parsedInstructionMeta.Source,
					}
				} else {
					if parsedInstructionMeta.Owner != "" {
						account = types.AccountIdentifier{
							Address: parsedInstructionMeta.Owner,
						}
					} else {
						if parsedInstructionMeta.Account != "" {
							account = types.AccountIdentifier{
								Address: parsedInstructionMeta.Account,
							}
						}
					}
				}

				operations = append(operations, &types.Operation{
					OperationIdentifier: &oi,
					Type:                opType,
					Account:             &account,
					Status:              &status,
					Metadata:            inInterface,
				})
			}
		}
	}
	return operations
}

func ToRosTxs(txs []ParsedTransactionWithMeta) []*RosettaTypes.Transaction {
	var rtxs []*RosettaTypes.Transaction
	for _, tx := range txs {
		rtx := ToRosTx(tx.Transaction)
		rtxs = append(rtxs, &rtx)
	}
	return rtxs
}
func ToRosTx(tx ParsedTransaction) RosettaTypes.Transaction {
	return RosettaTypes.Transaction{
		TransactionIdentifier: &RosettaTypes.TransactionIdentifier{
			Hash: tx.Signatures[0],
		},
		Operations: GetRosOperationsFromTx(tx, SuccessStatus),
		Metadata:   map[string]interface{}{},
	}
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
func EncodeBig(bigint *big.Int) string {
	nbits := bigint.BitLen()
	if nbits == 0 {
		return "0x0"
	}
	return fmt.Sprintf("%#x", bigint)
}

func convertTime(time uint64) int64 {
	return int64(time) * 1000
}

func GetWithNonce(m map[string]interface{}) (WithNonce, bool) {
	var withNonce WithNonce
	hasNonce := false
	if w, ok := m[WithNonceKey]; ok {
		j, _ := json.Marshal(w)
		json.Unmarshal(j, &withNonce)
		if len(withNonce.Account) > 0 {
			hasNonce = true
		}
	}
	return withNonce, hasNonce
}

func GetPriorityFee(m map[string]interface{}) PriorityFee {
	var priorityFee PriorityFee
	if w, ok := m[PriorityFeeKey]; ok {
		j, _ := json.Marshal(w)
		json.Unmarshal(j, &priorityFee)
	}
	return priorityFee
}
func GetTxFromStr(t string) (solPTypes.Transaction, error) {
	signedTx, err := base58.Decode(t)
	if err != nil {
		signedTx, err = hex.DecodeString(t)
		if err != nil {
			return solPTypes.Transaction{}, err
		}
	}

	tx, err := solPTypes.TransactionDeserialize(signedTx)
	if err != nil {
		return solPTypes.Transaction{}, err
	}

	return tx, nil
}
func ToParsedTransaction(tx solPTypes.Transaction) (ParsedTransaction, error) {
	ins := tx.Message.DecompileInstructions()
	var parsedIns []ParsedInstruction
	for _, v := range ins {
		p, err := ParseInstruction(v)
		if err != nil {
			//cannot parse
			p = ParsedInstruction{}
			return ParsedTransaction{}, err
		}
		parsedIns = append(parsedIns, p)
	}
	var acckeys []ParsedAccKey
	var sigs []string
	for _, v := range tx.Message.Accounts {
		acckeys = append(acckeys, ParsedAccKey{PubKey: v.ToBase58()})
	}
	for _, v := range tx.Signatures {
		sigs = append(sigs, base58.Encode(v[:]))
	}
	newTx := ParsedTransaction{
		Signatures: sigs,
		Message: ParsedMessage{
			Header:          tx.Message.Header,
			AccountKeys:     acckeys,
			RecentBlockhash: tx.Message.RecentBlockHash,
			Instructions:    parsedIns,
		},
	}
	return newTx, nil
}

type CreateAccountInstruction struct {
	Instruction Instruction
	Lamports    uint64
	Space       uint64
	Owner       common.PublicKey
}

func ParseInstruction(ins solPTypes.Instruction) (ParsedInstruction, error) {
	var parsedInstruction ParsedInstruction
	var err error

	switch ins.ProgramID {
	case common.SystemProgramID:
		parsedInstruction, err = ParseSystem(ins)
		break
	case common.TokenProgramID:
		parsedInstruction, err = ParseToken(ins)
		break
	case common.SPLAssociatedTokenAccountProgramID:
		parsedInstruction, err = ParseAssocToken(ins)
		break
	default:
		return parsedInstruction, fmt.Errorf("Cannot parse instruction")
	}
	if err != nil {
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

var (
	SystemProgramID                    = common.PublicKeyFromString("11111111111111111111111111111111")
	ConfigProgramID                    = common.PublicKeyFromString("Config1111111111111111111111111111111111111")
	StakeProgramID                     = common.PublicKeyFromString("Stake11111111111111111111111111111111111111")
	VoteProgramID                      = common.PublicKeyFromString("Vote111111111111111111111111111111111111111")
	BPFLoaderProgramID                 = common.PublicKeyFromString("BPFLoader1111111111111111111111111111111111")
	Secp256k1ProgramID                 = common.PublicKeyFromString("KeccakSecp256k11111111111111111111111111111")
	TokenProgramID                     = common.PublicKeyFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	SPLAssociatedTokenAccountProgramID = common.PublicKeyFromString("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")
)

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
	}
	return name
}

func ValueToBaseAmount(valueStr string) uint64 {
	var amount uint64
	valueStr = strings.Replace(valueStr, "-", "", -1)
	amt, err := strconv.ParseInt(valueStr, 10, 64)
	if err == nil {
		amount = uint64(amt)
	}
	return amount
}

type AssignInstruction struct {
	Instruction       Instruction
	AssignToProgramID common.PublicKey
}

type TransferInstruction struct {
	Instruction Instruction
	Lamports    uint64
}

type CreateAccountWithSeedInstruction struct {
	Instruction Instruction
	Base        common.PublicKey
	SeedLen     uint64
	Seed        string
	Lamports    uint64
	Space       uint64
	ProgramID   common.PublicKey
}

type AdvanceNonceAccountInstruction struct {
	Instruction Instruction
}
type WithdrawNonceAccountInstruction struct {
	Instruction Instruction
	Lamports    uint64
}

type InitializeNonceAccountInstruction struct {
	Instruction Instruction
	Auth        common.PublicKey
}
type AuthorizeNonceAccountInstruction struct {
	Instruction Instruction
	Auth        common.PublicKey
}
type AllocateInstruction struct {
	Instruction Instruction
	Space       uint64
}
type AllocateWithSeedInstruction struct {
	Instruction Instruction
	Base        common.PublicKey
	SeedLen     uint64
	Seed        string
	Space       uint64
	ProgramID   common.PublicKey
}
type AssignWithSeedInstruction struct {
	Instruction       Instruction
	Base              common.PublicKey
	SeedLen           uint64
	Seed              string
	AssignToProgramID common.PublicKey
}
type TransferWithSeedInstruction struct {
	Instruction Instruction
	Lamports    uint64
	SeedLen     uint64
	Seed        string
	ProgramID   common.PublicKey
}

func ParseSystem(ins solPTypes.Instruction) (ParsedInstruction, error) {
	var parsedInstruction ParsedInstruction
	var err error
	var s struct {
		Instruction InstructionInt
	}
	err = binstruct.UnmarshalLE(ins.Data, &s)
	var instructionType string
	var parsedInfo map[string]interface{}
	switch s.Instruction {
	case InstructionCreateAccount:
		var a CreateAccountInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "createAccount"
		parsedInfo = map[string]interface{}{
			"source":     ins.Accounts[0].PubKey.ToBase58(),
			"newAccount": ins.Accounts[1].PubKey.ToBase58(),
			"lamports":   a.Lamports,
			"space":      a.Space,
			"owner":      a.Owner.ToBase58(),
		}
		break
	case InstructionAssign:
		var a AssignInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "assign"
		parsedInfo = map[string]interface{}{
			"account": ins.Accounts[0].PubKey.ToBase58(),
			"owner":   a.AssignToProgramID.ToBase58(),
		}
		break
	case InstructionTransfer:
		var a TransferInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "transfer"
		parsedInfo = map[string]interface{}{
			"source":      ins.Accounts[0].PubKey.ToBase58(),
			"destination": ins.Accounts[1].PubKey.ToBase58(),
			"lamports":    a.Lamports,
		}
		parsedInstruction.Parsed = &InstructionInfo{}
		break
	case InstructionCreateAccountWithSeed:
		var a CreateAccountWithSeedInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "createAccountWithSeed"
		parsedInfo = map[string]interface{}{
			"source":     ins.Accounts[0].PubKey.ToBase58(),
			"newAccount": ins.Accounts[1].PubKey.ToBase58(),
			"base":       a.Base,
			"seed":       a.Seed,
			"space":      a.Space,
			"lamports":   a.Lamports,
			"owner":      a.ProgramID,
		}
		break
	case InstructionAdvanceNonceAccount:
		var a AdvanceNonceAccountInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "advanceNonce"
		parsedInfo = map[string]interface{}{
			"nonceAccount":            ins.Accounts[0].PubKey.ToBase58(),
			"recentBlockhashesSysvar": ins.Accounts[1].PubKey.ToBase58(),
			"nonceAuthority":          ins.Accounts[2].PubKey.ToBase58(),
		}
		break
	case InstructionWithdrawNonceAccount:
		var a WithdrawNonceAccountInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "withdrawFromNonce"
		parsedInfo = map[string]interface{}{
			"nonceAccount":            ins.Accounts[0].PubKey.ToBase58(),
			"destination":             ins.Accounts[1].PubKey.ToBase58(),
			"recentBlockhashesSysvar": ins.Accounts[2].PubKey.ToBase58(),
			"rentSysvar":              ins.Accounts[3].PubKey.ToBase58(),
			"nonceAuthority":          ins.Accounts[4].PubKey.ToBase58(),
			"lamports":                a.Lamports,
		}
		break
	case InstructionInitializeNonceAccount:
		var a InitializeNonceAccountInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "initializeNonce"
		parsedInfo = map[string]interface{}{
			"nonceAccount":            ins.Accounts[0].PubKey.ToBase58(),
			"recentBlockhashesSysvar": ins.Accounts[1].PubKey.ToBase58(),
			"rentSysvar":              ins.Accounts[2].PubKey.ToBase58(),
			"nonceAuthority":          a.Auth.ToBase58(),
		}
		break
	case InstructionAuthorizeNonceAccount:
		var a AuthorizeNonceAccountInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "authorizeNonce"
		parsedInfo = map[string]interface{}{
			"nonceAccount":   ins.Accounts[0].PubKey.ToBase58(),
			"nonceAuthority": ins.Accounts[1].PubKey.ToBase58(),
			"rentSysvar":     ins.Accounts[2].PubKey.ToBase58(),
			"newAuthorized":  a.Auth.ToBase58(),
		}
		break
	case InstructionAllocate:
		var a AllocateInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "allocate"
		parsedInfo = map[string]interface{}{
			"account": ins.Accounts[0].PubKey.ToBase58(),
			"space":   a.Space,
		}
		break
	case InstructionAllocateWithSeed:
		var a AllocateWithSeedInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "allocateWithSeed"
		parsedInfo = map[string]interface{}{
			"account": ins.Accounts[0].PubKey.ToBase58(),
			"base":    a.Base,
			"seed":    a.Seed,
			"space":   a.Space,
			"owner":   a.ProgramID.ToBase58(),
		}
		break
	case InstructionAssignWithSeed:
		var a AssignWithSeedInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "assignWithSeed"
		parsedInfo = map[string]interface{}{
			"account": ins.Accounts[0].PubKey.ToBase58(),
			"base":    a.Base,
			"seed":    a.Seed,
			"owner":   a.AssignToProgramID.ToBase58(),
		}
		break
	case InstructionTransferWithSeed:
		var a TransferWithSeedInstruction
		err = binstruct.UnmarshalLE(ins.Data, &a)
		instructionType = "transferWithSeed"
		parsedInfo = map[string]interface{}{
			"source":      ins.Accounts[0].PubKey.ToBase58(),
			"sourceBase":  ins.Accounts[1].PubKey.ToBase58(),
			"destination": ins.Accounts[2].PubKey.ToBase58(),
			"lamports":    a.Lamports,
			"sourceSeed":  a.Seed,
			"sourceOwner": a.ProgramID.ToBase58(),
		}
		break
	}
	parsedInstruction.Parsed = &InstructionInfo{
		Info:            parsedInfo,
		InstructionType: instructionType,
	}
	return parsedInstruction, err
}
