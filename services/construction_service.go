// Copyright 2020 Coinbase, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package services

import (
	"context"
	"fmt"
	"github.com/blocto/solana-go-sdk/program/system"
	"github.com/imerkle/rosetta-solana-go/solana/parse"
	stypes "github.com/imerkle/rosetta-solana-go/solana/shared_types"
	"log"
	"strconv"
	"strings"

	"github.com/blocto/solana-go-sdk/common"
	solPTypes "github.com/blocto/solana-go-sdk/types"
	"github.com/imerkle/rosetta-solana-go/configuration"
	solanago "github.com/imerkle/rosetta-solana-go/solana"
	"github.com/imerkle/rosetta-solana-go/solana/operations"
	"github.com/mitchellh/copystructure"
	"github.com/mr-tron/base58"

	"github.com/coinbase/rosetta-sdk-go/types"
)

// ConstructionAPIService implements the server.ConstructionAPIServicer interface.
type ConstructionAPIService struct {
	config       *configuration.Configuration
	client       *solanago.Client
	directClient *solanago.DirectClient
}

// NewConstructionAPIService creates a new instance of a ConstructionAPIService.
func NewConstructionAPIService(
	cfg *configuration.Configuration,
	client *solanago.Client,
) *ConstructionAPIService {
	var directClient = solanago.NewDirectClient(cfg.GethURL)
	return &ConstructionAPIService{
		config:       cfg,
		client:       client,
		directClient: directClient,
	}
}

// ConstructionDerive implements the /construction/derive endpoint.
func (s *ConstructionAPIService) ConstructionDerive(
	ctx context.Context,
	request *types.ConstructionDeriveRequest,
) (*types.ConstructionDeriveResponse, *types.Error) {
	addr := base58.Encode(request.PublicKey.Bytes)
	return &types.ConstructionDeriveResponse{
		AccountIdentifier: &types.AccountIdentifier{
			Address: addr,
		},
	}, nil
}

// ConstructionPreprocess implements the /construction/preprocess
// endpoint.
func (s *ConstructionAPIService) ConstructionPreprocess(
	ctx context.Context,
	request *types.ConstructionPreprocessRequest,
) (*types.ConstructionPreprocessResponse, *types.Error) {
	log.Printf("START /construction/preprocess")
	log.Printf("request.Metadata=%+v\n", request.Metadata)

	withNonce, hasNonce := solanago.GetWithNonce(request.Metadata)
	log.Printf("withNonce=%+v\n", withNonce)
	if hasNonce {
		log.Printf("inside hasNonce=true")
		acc, _ := s.directClient.GetAccountInfoParsed(ctx, withNonce.Account)
		withNonce.Authority = acc.Data.Parsed.Info.Authority
	} else {
		log.Printf("inside hasNonce=false")
	}

	priorityFee := solanago.GetPriorityFee(request.Metadata)
	log.Printf("priorityFee=%+v\n", priorityFee)

	log.Printf("request.Operations=%+v\n", request.Operations)

	var matchedOperationHashMap = make(map[int64]bool)

	var SplSystemAccMap = make(map[int64]stypes.SplAccounts)
	for _, op := range request.Operations {
		LogOperation(op)

		var cont bool
		var matched *types.Operation
		cont, matched = FindMatch(request.Operations, op, matchedOperationHashMap)
		if cont {
			continue
		}
		if matched != nil && op.Type == stypes.SplToken__TransferWithSystem {
			SplSystemAccMap[op.OperationIdentifier.Index] = stypes.SplAccounts{
				Source:      op.Account.Address,
				Destination: matched.Account.Address,
				Mint:        op.Amount.Currency.Symbol,
			}
			matchedOperationHashMap[op.OperationIdentifier.Index] = true
		}
	}

	var constructionMetaData = ConstructionMetadata{
		PriorityFee: priorityFee,
		WithNonce:   withNonce,
	}

	_, instructions, _ := ToInstructions(request.Operations, constructionMetaData)

	instructions = AdvanceNonce(withNonce, instructions)
	signers := GetUniqueSigners(instructions)

	var feeCalculation = stypes.FeeCalculation{
		NumberOfInstructions: strconv.Itoa(len(instructions)),
		NumberOfSigners:      strconv.Itoa(len(signers)),
	}
	log.Printf("instructions.len=%+v\n", len(instructions))

	log.Printf("END /construction/preprocess")
	return &types.ConstructionPreprocessResponse{
		Options: map[string]interface{}{
			stypes.WithNonceKey:       withNonce,
			stypes.FeeCalculationKey:  feeCalculation,
			stypes.PriorityFeeKey:     priorityFee,
			stypes.SplSystemAccMapKey: SplSystemAccMap,
		},
	}, nil
}

// ConstructionHash implements the /construction/hash endpoint.
func (s *ConstructionAPIService) ConstructionHash(
	ctx context.Context,
	request *types.ConstructionHashRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	log.Printf("START /construction/hash")

	tx, err := solanago.GetTxFromStr(request.SignedTransaction)
	if err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}
	hash := base58.Encode(tx.Signatures[0])

	log.Printf("hash=%s\n", hash)
	log.Printf("END /construction/hash")
	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: &types.TransactionIdentifier{
			Hash: hash,
		},
	}, nil
}

// ConstructionMetadata implements the /construction/metadata endpoint.
func (s *ConstructionAPIService) ConstructionMetadata(
	ctx context.Context,
	request *types.ConstructionMetadataRequest,
) (*types.ConstructionMetadataResponse, *types.Error) {
	log.Printf("START /construction/metadata")
	if s.config.Mode != configuration.Online {
		return nil, ErrUnavailableOffline
	}

	var hash string
	var blockNumber uint64
	var feeCalculator stypes.FeeCalculator
	withNonce, hasNonce := solanago.GetWithNonce(request.Options)
	log.Printf("withNonce=%s\n", withNonce)
	log.Printf("hasNonce=%t\n", hasNonce)

	priorityFee := solanago.GetPriorityFee(request.Options)
	log.Printf("priorityFee=%+v\n", priorityFee)

	feeCalculation := solanago.GetFeeCalculation(request.Options)

	if hasNonce {
		log.Printf("inside hasNonce=true")
		acc, _ := s.directClient.GetAccountInfoParsed(ctx, withNonce.Account)
		//s.client.Rpc.GetAccountInfo(ctx, withNonce.Account)
		withNonce.Authority = acc.Data.Parsed.Info.Authority
		hash = acc.Data.Parsed.Info.BlockHash
		var ssFeeCalculator = acc.Data.Parsed.Info.FeeCalculator
		feeCalculator = stypes.FeeCalculator{LamportsPerSignature: solanago.ValueToBaseAmount(ssFeeCalculator.LamportsPerSignature)}
		log.Printf("acc=%+v\n", acc)
	} else {
		log.Printf("inside hasNonce=false")
		status, _, _, _, err := s.client.Status(ctx)
		if err != nil {
			return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
		}
		hash = status.Hash
		blockNumber = uint64(status.Index)
		// hardcoded to 5000 as the fee can not be fetched anymore without first building a transaction
		feeCalculator = stypes.FeeCalculator{LamportsPerSignature: 5000}
		log.Printf("feeCalculator=%d\n", feeCalculator)
		log.Printf("blockHash=%s\n", hash)
	}

	var SplTokenAccMap = make(map[string]stypes.SplAccounts)

	if w, ok := request.Options[stypes.SplSystemAccMapKey]; ok {
		w1 := w.(map[string]interface{})
		if err := unmarshalJSONMap(w1, &SplTokenAccMap); err != nil {
			return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
		}

		for k, v := range SplTokenAccMap {

			source, _ := s.directClient.GetTokenAccountByOwner(ctx, v.Source)
			destination, _ := s.directClient.GetTokenAccountByMint(ctx, v.Destination, v.Mint)
			SplTokenAccMap[k] = stypes.SplAccounts{
				Source:      source,
				Destination: destination,
				Mint:        v.Mint,
			}
		}
	}

	meta, _ := marshalJSONMap(ConstructionMetadata{
		BlockHash:         hash,
		BlockNumber:       blockNumber,
		PriorityFee:       priorityFee,
		FeeCalculator:     feeCalculator,
		SplTokenAccMapKey: SplTokenAccMap,
		WithNonce:         withNonce,
		FeeCalculation:    feeCalculation,
	})

	log.Printf("meta=%+v\n", meta)
	log.Printf("END /construction/metadata")

	return &types.ConstructionMetadataResponse{
		Metadata: meta,
		SuggestedFee: []*types.Amount{
			{
				Value:    strconv.FormatInt(int64(feeCalculator.LamportsPerSignature), 10),
				Currency: stypes.Currency,
			},
		},
	}, nil
}

// ConstructionPayloads implements the /construction/payloads endpoint.
func (s *ConstructionAPIService) ConstructionPayloads(
	ctx context.Context,
	request *types.ConstructionPayloadsRequest,
) (*types.ConstructionPayloadsResponse, *types.Error) {
	log.Printf("START /construction/payloads")

	// Convert map to Metadata struct
	var meta ConstructionMetadata

	if err := unmarshalJSONMap(request.Metadata, &meta); err != nil {
		log.Printf("err=%s", err)
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}

	log.Printf("meta=%+v\n", meta)
	ops := request.Operations

	feePayer, instructions, toInstructionsErr := ToInstructions(ops, meta)
	if toInstructionsErr != nil {
		return nil, toInstructionsErr
	}
	// this list is without the nonce-advance and so we can use the first signer as the default if needed
	signers := GetUniqueSigners(instructions)

	if feePayer == (common.PublicKey{}) {
		feePayer = common.PublicKeyFromString(signers[0])
	}

	blockHash := meta.BlockHash
	var message solPTypes.Message

	instructions = AdvanceNonce(meta.WithNonce, instructions)
	_, hasNonce := solanago.GetWithNonce(request.Metadata)
	if hasNonce {
		message = solPTypes.NewMessage(solPTypes.NewMessageParam{FeePayer: feePayer, Instructions: instructions, RecentBlockhash: ""})
	} else {
		message = solPTypes.NewMessage(solPTypes.NewMessageParam{FeePayer: feePayer, Instructions: instructions, RecentBlockhash: blockHash})
	}

	//unsigned signature
	var sig []solPTypes.Signature
	x := make([]byte, 64)
	for i := 0; i < int(message.Header.NumRequireSignatures); i++ {
		sig = append(sig, x)
	}
	tx := solPTypes.Transaction{
		Signatures: sig,
		Message:    message,
	}
	tx.Message.RecentBlockHash = blockHash
	msgBytes, _ := tx.Message.Serialize()
	var signingPayloads []*types.SigningPayload
	for _, sg := range signers {
		log.Printf("signer=%+v", sg)
		signingPayloads = append(signingPayloads, &types.SigningPayload{
			AccountIdentifier: &types.AccountIdentifier{
				Address: sg,
			},
			Bytes:         msgBytes,
			SignatureType: types.Ed25519,
		})
	}

	txUnsigned, err := tx.Serialize()

	if err != nil {
		log.Printf("err=%s", err)
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}

	log.Printf("END /construction/payloads")
	return &types.ConstructionPayloadsResponse{
		UnsignedTransaction: base58.Encode(txUnsigned),
		Payloads:            signingPayloads,
	}, nil
}

func AdvanceNonce(withNonce stypes.WithNonce, instructions []solPTypes.Instruction) []solPTypes.Instruction {
	// If a nonce is specified we have to advance it
	if (withNonce != stypes.WithNonce{}) {
		log.Printf("ToInstructions withNonce=%+v\n", withNonce)
		ins := system.AdvanceNonceAccount(system.AdvanceNonceAccountParam{p(withNonce.Account), p(withNonce.Authority)})
		instructions = append([]solPTypes.Instruction{ins}, instructions...)
	}
	return instructions
}

// ConstructionCombine implements the /construction/combine
// endpoint.
func (s *ConstructionAPIService) ConstructionCombine(
	ctx context.Context,
	request *types.ConstructionCombineRequest,
) (*types.ConstructionCombineResponse, *types.Error) {
	log.Printf("START /construction/combine")

	tx, err := solanago.GetTxFromStr(request.UnsignedTransaction)
	if err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}
	var pubKeys []common.PublicKey
	for _, s := range request.Signatures {
		pubKeys = append(pubKeys, common.PublicKeyFromBytes(s.PublicKey.Bytes))
	}
	positions, errr := GetSigningKeypairPositions(tx.Message, pubKeys)
	if errr != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}
	for i, p := range positions {
		log.Printf("p=%d, i=%d", p, i)
		tx.Signatures[p] = request.Signatures[i].Bytes
		log.Printf("tx.Signatures[p]=%s", base58.Encode(tx.Signatures[p]))
	}
	signedTx, err := tx.Serialize()
	if err != nil {
		return nil, wrapErr(ErrSignatureInvalid, err)
	}

	log.Printf("END /construction/combine")
	return &types.ConstructionCombineResponse{
		SignedTransaction: base58.Encode(signedTx),
	}, nil
}

// ConstructionParse implements the /construction/parse endpoint.
func (s *ConstructionAPIService) ConstructionParse(
	ctx context.Context,
	request *types.ConstructionParseRequest,
) (*types.ConstructionParseResponse, *types.Error) {

	tx, err := solanago.GetTxFromStr(request.Transaction)
	if err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}

	var signers []*types.AccountIdentifier
	sgns := GetUniqueSigners(tx.Message.DecompileInstructions())
	for _, v := range sgns {
		signers = append(signers, &types.AccountIdentifier{
			Address: v,
		})
	}
	parsedTx, err := parse.ToParsedTransaction(tx)
	if err != nil {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, err)
	}

	operations := solanago.GetRosOperationsFromTx(parsedTx, "")

	resp := &types.ConstructionParseResponse{
		Operations:               operations,
		AccountIdentifierSigners: signers,
	}
	return resp, nil
}

// ConstructionSubmit implements the /construction/submit endpoint.
func (s *ConstructionAPIService) ConstructionSubmit(
	ctx context.Context,
	request *types.ConstructionSubmitRequest,
) (*types.TransactionIdentifierResponse, *types.Error) {
	log.Printf("START /construction/submit")

	version, _ := s.client.Rpc.GetVersion(ctx)
	log.Printf("s.config.GethURL=%s\n", s.config.GethURL)
	log.Printf("s.config.Mode=%v\n", s.config.Mode)
	log.Printf("version.FeatureSet=%d\n", version.FeatureSet)
	log.Printf("version.SolanaCore=%s\n", version.SolanaCore)
	log.Printf("request.SignedTransaction\n%s\n", request.SignedTransaction)

	log.Printf("before Rpc.SendTransaction")
	if s.config.Mode != configuration.Online {
		return nil, ErrUnavailableOffline
	}
	decode, err2 := base58.Decode(request.SignedTransaction)
	if err2 != nil {
		log.Printf("err=%s", err2)
		return nil, wrapErr(ErrBroadcastFailed, err2)
	}
	transaction, err := solPTypes.TransactionDeserialize(decode)
	if err != nil {
		log.Printf("err=%s", err)
		return nil, wrapErr(ErrBroadcastFailed, err)
	}

	hash, err := s.client.Rpc.SendTransaction(ctx, transaction)
	if err != nil {
		return nil, wrapErr(ErrBroadcastFailed, err)
	}
	log.Printf("after Rpc.SendTransaction")
	log.Printf("hash=%s\n", hash)

	txIdentifier := &types.TransactionIdentifier{
		Hash: hash,
	}

	log.Printf("END /construction/submit")
	return &types.TransactionIdentifierResponse{
		TransactionIdentifier: txIdentifier,
	}, nil
}

func LogOperation(op *types.Operation) {
	log.Printf("op=%+v\n", op)
	log.Printf("op.Type=%v\n", op.Type)
	if op.Account != nil {
		log.Printf("op.Account.Address=%s\n", op.Account.Address)
	}
	if op.Amount != nil {
		log.Printf("op.Amount.Value=%s\n", op.Amount.Value)
	}
	if op.Metadata != nil {
		log.Printf("op.Metadata=%s\n", op.Metadata)
	}
}

func FindMatch(ops []*types.Operation, op *types.Operation, matchedOperationHashMap map[int64]bool) (bool, *types.Operation) {
	if _, ok := matchedOperationHashMap[op.OperationIdentifier.Index]; ok {
		return true, nil
	}
	var matched *types.Operation = nil
	for _, v := range ops {
		if op.OperationIdentifier.Index == v.OperationIdentifier.Index {
			continue
		}
		if _, ok := matchedOperationHashMap[v.OperationIdentifier.Index]; ok {
			continue
		}
		if v.Type != op.Type {
			continue
		}
		if v.Amount != nil {
			if v.Amount.Currency.Symbol != op.Amount.Currency.Symbol {
				continue
			}
			if solanago.ValueToBaseAmount(v.Amount.Value) != solanago.ValueToBaseAmount(op.Amount.Value) {
				continue
			} else {
				opisNegative := strings.Contains(op.Amount.Value, "-")
				visNegative := strings.Contains(v.Amount.Value, "-")
				if (opisNegative && visNegative) || (!opisNegative && !visNegative) {
					continue
				}
			}
		}
		return false, v
	}
	return false, matched
}

func NewMessageWithNonce(feePayer common.PublicKey, instructions []solPTypes.Instruction, nonceAccountPubkey common.PublicKey, nonceAuthorityPubkey common.PublicKey) solPTypes.Message {
	//ins := system.AdvanceNonceAccount(system.AdvanceNonceAccountParam{nonceAccountPubkey, nonceAuthorityPubkey})
	//instructions = append([]solPTypes.Instruction{ins}, instructions...)
	message := solPTypes.NewMessage(solPTypes.NewMessageParam{FeePayer: feePayer, Instructions: instructions, RecentBlockhash: ""})
	return message
}

func GetSigningKeypairPositions(message solPTypes.Message, pubKeys []common.PublicKey) ([]uint, *types.Error) {
	if len(message.Accounts) < int(message.Header.NumRequireSignatures) {
		return nil, wrapErr(ErrUnableToParseIntermediateResult, fmt.Errorf("invalid positions"))
	}
	signedKeys := message.Accounts[0:message.Header.NumRequireSignatures]
	var positions []uint
	for _, p := range pubKeys {
		index := indexOf(p, signedKeys)
		if index > -1 {
			positions = append(positions, uint(index))
		}
	}
	return positions, nil
}
func indexOf(element common.PublicKey, data []common.PublicKey) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func GetUniqueSigners(ins []solPTypes.Instruction) []string {
	var signers []string
	var signersMap = make(map[string]bool)
	for _, v := range ins {
		for _, v1 := range v.Accounts {
			address := v1.PubKey.ToBase58()
			if v1.IsSigner {
				if _, ok := signersMap[address]; !ok {
					signersMap[address] = true
					signers = append(signers, address)
				}
			}
		}
	}
	return signers
}

func ToInstructions(ops []*types.Operation, meta ConstructionMetadata) (common.PublicKey, []solPTypes.Instruction, *types.Error) {
	log.Printf("START ToInstructions")
	var instructions []solPTypes.Instruction
	var feePayer common.PublicKey
	var matchedOperationHashMap = make(map[int64]bool)

	for _, op := range ops {
		LogOperation(op)
		var cont bool
		var matched *types.Operation
		cont, matched = FindMatch(ops, op, matchedOperationHashMap)
		if cont {
			log.Printf("op is alread in matchedOperationHashMap -> continue with next entry")
			continue
		}

		if matched == nil && op.Amount != nil {
			return common.PublicKey{}, nil, wrapErr(ErrUnableToParseIntermediateResult, fmt.Errorf("Invalid Operation Request. Please check format"))
		}

		opCopy, err := copystructure.Copy(*op)
		if err != nil {
			return common.PublicKey{}, nil, wrapErr(ErrUnclearIntent, fmt.Errorf("Cannot deep copy operations"))

		}
		tp := opCopy.(types.Operation)
		tmpOP := &tp

		if tmpOP.Metadata == nil {
			tmpOP.Metadata = make(map[string]interface{})
		}
		if matched != nil {
			fromOp := tmpOP
			fromAdd := fromOp.Account.Address
			toOp := matched
			toAdd := toOp.Account.Address
			tmpOP.Account = fromOp.Account
			tmpOP.Metadata["source"] = fromAdd
			tmpOP.Metadata["destination"] = toAdd
			tmpOP.Amount = toOp.Amount

			matchedOperationHashMap[fromOp.OperationIdentifier.Index] = true
			matchedOperationHashMap[toOp.OperationIdentifier.Index] = true

		} else {
			matchedOperationHashMap[op.OperationIdentifier.Index] = true
		}

		log.Printf("tmpOP.Type=%s\n", tmpOP.Type)
		switch strings.Split(tmpOP.Type, stypes.Separator)[0] {
		case "System":
			s := operations.SystemOperationMetadata{}
			s.SetMeta(tmpOP, meta.PriorityFee)
			instructions = append(instructions, s.ToInstructions(tmpOP.Type)...)
			break
		case "SplToken":
			s := operations.SplTokenOperationMetadata{}
			s.SetMeta(tmpOP, meta.SplTokenAccMapKey)
			instructions = append(instructions, s.ToInstructions(tmpOP.Type)...)
			break
		case "SplAssociatedTokenAccount":
			s := operations.SplAssociatedTokenAccountOperationMetadata{}
			s.SetMeta(tmpOP)
			instructions = append(instructions, s.ToInstructions(tmpOP.Type)...)
			break
		case "Stake":
			s := operations.StakeOperationMetadata{}
			s.SetMeta(tmpOP, meta.PriorityFee)
			instructions = append(instructions, s.ToInstructions(tmpOP.Type)...)
			if tmpOP.Type == stypes.Stake__WithdrawStake && s.FeePayer != "" {
				feePayer = common.PublicKeyFromString(s.FeePayer)
			}
			break
		default:
			return common.PublicKey{}, nil, wrapErr(ErrUnableToParseIntermediateResult, fmt.Errorf("Operation not implemented for construction"))
		}
	}

	log.Printf("There are in total %v instructions", len(instructions))
	for i, in := range instructions {
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
	return feePayer, instructions, nil
}

func p(a string) common.PublicKey {
	return common.PublicKeyFromString(a)
}
