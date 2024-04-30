package solanago

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	stypes "github.com/imerkle/rosetta-solana-go/solana/shared_types"
	"math/big"
	"strconv"
	"strings"

	solPTypes "github.com/blocto/solana-go-sdk/types"
	"github.com/coinbase/rosetta-sdk-go/types"
	RosettaTypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/mr-tron/base58"

	"github.com/iancoleman/strcase"
)

func IsBalanceChanging(opType string) bool {
	a := false
	switch opType {
	case stypes.System__CreateAccount, stypes.System__WithdrawFromNonce, stypes.System__Transfer, stypes.SplToken__Transfer, stypes.SplToken__TransferChecked, "Stake__Split", "Stake__Withdraw", "Vote__Withdraw", stypes.SplToken__TransferNew, stypes.SplToken__TransferWithSystem:
		a = true
	}
	return a
}

func getOperationTypeWithProgram(program string, s string) string {
	toPascal := strcase.ToCamel(program)

	newStr := fmt.Sprint(
		toPascal,
		stypes.Separator,
		strcase.ToCamel(s),
	)
	return newStr
}
func getOperationType(s string) string {
	x := strings.Split(s, stypes.Separator)
	if len(x) < 2 {
		return stypes.Unknown
	}
	return getOperationTypeWithProgram(x[0], x[1])
}
func split_at(at int, input []byte) ([]byte, []byte) {
	return input[0:1], input[1:]
}
func GetRosOperationsFromTx(tx stypes.ParsedTransaction, status string) []*types.Operation {
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
				Type:                stypes.Unknown,
				Status:              &status,
				Metadata:            inInterface,
			})
		} else {

			jsonString, _ := json.Marshal(ins.Parsed.Info)

			parsedInstructionMeta := stypes.ParsedInstructionMeta{}
			var parsedInstructionMetaInterface interface{}
			json.Unmarshal(jsonString, &parsedInstructionMeta)
			json.Unmarshal(jsonString, &parsedInstructionMetaInterface)

			var inInterface map[string]interface{}
			inrec, _ := json.Marshal(parsedInstructionMetaInterface)
			json.Unmarshal(inrec, &inInterface)

			opType := getOperationTypeWithProgram(ins.Program, ins.Parsed.InstructionType)
			if !Contains(stypes.OperationTypes, opType) {
				inInterface["instruction_type"] = ins.Parsed.InstructionType
				inInterface["program"] = ins.Program
				opType = "Unknown"
			}
			if IsBalanceChanging(opType) {
				if parsedInstructionMeta.Decimals == 0 {
					parsedInstructionMeta.Decimals = stypes.Decimals
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
							Symbol:   stypes.Symbol,
							Decimals: stypes.Decimals,
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

func ToRosTxs(txs []stypes.ParsedTransactionWithMeta) []*RosettaTypes.Transaction {
	var rtxs []*RosettaTypes.Transaction
	for _, tx := range txs {
		rtx := ToRosTx(tx.Transaction)
		rtxs = append(rtxs, &rtx)
	}
	return rtxs
}
func ToRosTx(tx stypes.ParsedTransaction) RosettaTypes.Transaction {
	return RosettaTypes.Transaction{
		TransactionIdentifier: &RosettaTypes.TransactionIdentifier{
			Hash: tx.Signatures[0],
		},
		Operations: GetRosOperationsFromTx(tx, stypes.SuccessStatus),
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

func GetWithNonce(m map[string]interface{}) (stypes.WithNonce, bool) {
	var withNonce stypes.WithNonce
	hasNonce := false
	if w, ok := m[stypes.WithNonceKey]; ok {
		j, _ := json.Marshal(w)
		json.Unmarshal(j, &withNonce)
		if len(withNonce.Account) > 0 {
			hasNonce = true
		}
	}
	return withNonce, hasNonce
}

func GetPriorityFee(m map[string]interface{}) stypes.PriorityFee {
	var priorityFee stypes.PriorityFee
	if w, ok := m[stypes.PriorityFeeKey]; ok {
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

func ValueToBaseAmount(valueStr string) uint64 {
	var amount uint64
	valueStr = strings.Replace(valueStr, "-", "", -1)
	amt, err := strconv.ParseInt(valueStr, 10, 64)
	if err == nil {
		amount = uint64(amt)
	}
	return amount
}
