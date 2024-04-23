package solanago

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	solPTypes "github.com/blocto/solana-go-sdk/types"
	"io/ioutil"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Context struct {
	Slot uint64 `json:"slot"`
}

type GeneralResponse struct {
	JsonRPC string        `json:"jsonrpc"`
	ID      uint64        `json:"id"`
	Error   ErrorResponse `json:"error"`
}

type DirectClient struct {
	endpoint string
}

func NewDirectClient(endpoint string) *DirectClient {
	return &DirectClient{endpoint: endpoint}
}

func (s *DirectClient) request(ctx context.Context, method string, params []interface{}, response interface{}) error {
	// post data
	j, err := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      0,
		"method":  method,
		"params":  params,
	})
	if err != nil {
		return err
	}

	// post request
	req, err := http.NewRequestWithContext(ctx, "POST", s.endpoint, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	// http client and send request
	httpclient := &http.Client{}
	res, err := httpclient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// parse body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(body) != 0 {
		if err := json.Unmarshal(body, &response); err != nil {
			return err
		}
	}

	// return result
	if res.StatusCode < 200 || res.StatusCode > 300 {
		return fmt.Errorf("get status code: %d", res.StatusCode)
	}
	return nil
}

func (s *DirectClient) CallRequest(ctx context.Context, method string, params []interface{}) (interface{}, error) {
	var out interface{}
	var err error

	if err != nil {
		return nil, err
	}
	return out, nil
}

type GetRecentBlockHashResponseValue struct {
	Blockhash     string        `json:"blockhash"`
	FeeCalculator FeeCalculator `json:"feeCalculator"`
}

type GetRecentBlockHashResponse struct {
	Context Context                         `json:"context"`
	Value   GetRecentBlockHashResponseValue `json:"value"`
}

func (s *DirectClient) GetRecentBlockhash(ctx context.Context) (GetRecentBlockHashResponse, error) {
	var res = struct {
		GeneralResponse
		Result struct {
			Context Context                         `json:"context"`
			Value   GetRecentBlockHashResponseValue `json:"value"`
		} `json:"result"`
	}{}
	err := s.request(ctx, "getRecentBlockhash", []interface{}{}, &res)
	if err != nil {
		return GetRecentBlockHashResponse{}, err
	}
	if res.Error != (ErrorResponse{}) {
		return GetRecentBlockHashResponse{}, errors.New(res.Error.Message)
	}
	result := GetRecentBlockHashResponse{
		Context: res.Result.Context,
		Value:   res.Result.Value,
	}
	return result, nil
}

type GetAccountInfoParsedResponse struct {
	Lamports  uint64  `json:"lamports"`
	Owner     string  `json:"owner"`
	Excutable bool    `json:"excutable"`
	RentEpoch uint64  `json:"rentEpoch"`
	Data      AccData `json:"data"`
}

type Info struct {
	Authority     string              `json:"authority"`
	BlockHash     string              `json:"blockhash"`
	FeeCalculator FeeCalculatorString `json:"feeCalculator"`
}
type Parsed struct {
	Info Info `json:"info"`
}
type AccData struct {
	Parsed Parsed `json:"parsed"`
}

type FeeCalculator struct {
	LamportsPerSignature uint64 `json:"lamportsPerSignature"`
}

type FeeCalculatorString struct {
	LamportsPerSignature string `json:"lamportsPerSignature"`
}

type ParsedTransaction struct {
	Signatures []string      `json:"signatures"`
	Message    ParsedMessage `json:"message"`
}

type ParsedMessage struct {
	Header          solPTypes.MessageHeader `json:"header"`
	AccountKeys     []ParsedAccKey          `json:"accountKeys"`
	RecentBlockhash string                  `json:"recentBlockhash"`
	Instructions    []ParsedInstruction     `json:"instructions"`
}

type ParsedAccKey struct {
	PubKey     string `json:"pubkey"`
	IsSigner   bool   `json:"signer"`
	IsWritable bool   `json:"writable"`
}

type ParsedInstruction struct {
	Accounts  []string         `json:"accounts,omitempty"`
	Data      string           `json:"data,omitempty"`
	Parsed    *InstructionInfo `json:"parsed,omitempty"`
	Program   string           `json:"program,omitempty"`
	ProgramID string           `json:"programId"`
}

type InstructionInfo struct {
	Info            map[string]interface{} `json:"info"`
	InstructionType string                 `json:"type"`
}

type Reward struct {
	Pubkey      string `json:"pubkey"`
	Lamports    int64  `json:"lamports"`
	PostBalance uint64 `json:"postBalance"`
	RewardType  string `json:"rewardType"` // type of reward: "fee", "rent", "voting", "staking"
}

type GetConfirmBlockParsedResponse struct {
	Blockhash         string                      `json:"blockhash"`
	PreviousBlockhash string                      `json:"previousBlockhash"`
	ParentSlot        uint64                      `json:"parentSlot"`
	BlockTime         int64                       `json:"blockTime"`
	Transactions      []ParsedTransactionWithMeta `json:"transactions"`
	Rewards           []Reward                    `json:"rewards"`
}

type ParsedTransactionWithMeta struct {
	Meta        TransactionMeta   `json:"meta"`
	Transaction ParsedTransaction `json:"transaction"`
}

type Transaction struct {
	Signatures []string `json:"signatures"`
	Message    Message  `json:"message"`
}

type Message struct {
	Header          MessageHeader `json:"header"`
	AccountKeys     []string      `json:"accountKeys"`
	RecentBlockhash string        `json:"recentBlockhash"`
	Instructions    []Instruction `json:"instructions"`
}

type Instruction struct {
	ProgramIDIndex uint64   `json:"programIdIndex"`
	Accounts       []uint64 `json:"accounts"`
	Data           string   `json:"data"`
}

type MessageHeader struct {
	NumRequiredSignatures       uint8 `json:"numRequiredSignatures"`
	NumReadonlySignedAccounts   uint8 `json:"numReadonlySignedAccounts"`
	NumReadonlyUnsignedAccounts uint8 `json:"numReadonlyUnsignedAccounts"`
}

type TransactionMeta struct {
	Fee               uint64   `json:"fee"`
	PreBalances       []int64  `json:"preBalances"`
	PostBalances      []int64  `json:"postBalances"`
	LogMessages       []string `json:"logMesssages"`
	InnerInstructions []struct {
		Index        uint64        `json:"index"`
		Instructions []Instruction `json:"instructions"`
	} `json:"innerInstructions"`
	Err    interface{}            `json:"err"`
	Status map[string]interface{} `json:"status"`
}

type GetConfirmedTransactionResponse struct {
	Slot        uint64          `json:"slot"`
	Meta        TransactionMeta `json:"meta"`
	Transaction Transaction     `json:"transaction"`
}

type GetConfirmedTransactionParsedResponse struct {
	Slot        uint64            `json:"slot"`
	Meta        TransactionMeta   `json:"meta"`
	Transaction ParsedTransaction `json:"transaction"`
}

type TokenAmount struct {
	Amount         string  `json:"amount"`
	Decimals       int32   `json:"decimals"`
	UIAmount       float64 `json:"uiAmount"`
	UIAmountString string  `json:"uiAmountString"`
}
type ParsedAccountInfo struct {
	Delegate        string      `json:"delegate"`
	DelegatedAmount TokenAmount `json:"delegatedAmount,omitempty"`
	IsInitialized   bool        `json:"isInitialized"`
	IsNative        bool        `json:"isNative"`
	Mint            string      `json:"mint"`
	Owner           string      `json:"owner"`
	TokenAmount     TokenAmount `json:"tokenAmount"`
}
type ParsedAccountData struct {
	AccountType string            `json:"accountType"`
	Info        ParsedAccountInfo `json:"info"`
}
type Data struct {
	Parsed  ParsedAccountData `json:"parsed"`
	Program string            `json:"program"`
}
type Account struct {
	Data       Data   `json:"data"`
	Executable bool   `json:"executable"`
	Lamports   int64  `json:"lamports"`
	Owner      string `json:"owner"`
	RentEpoch  int64  `json:"rentEpoch"`
}
type Accounts struct {
	Account Account `json:"account"`
	Pubkey  string  `json:"pubkey,omitempty"`
}

func (s *DirectClient) GetAccountInfoParsed(ctx context.Context, account string) (GetAccountInfoParsedResponse, error) {
	res := struct {
		GeneralResponse
		Result struct {
			Context Context                      `json:"context"`
			Value   GetAccountInfoParsedResponse `json:"value"`
		} `json:"result"`
	}{}
	err := s.request(ctx, "getAccountInfo", []interface{}{account, map[string]interface{}{"encoding": "jsonParsed"}}, &res)
	if err != nil {
		return GetAccountInfoParsedResponse{}, err
	}
	if res.Error != (ErrorResponse{}) {
		return GetAccountInfoParsedResponse{}, errors.New(res.Error.Message)
	}
	return res.Result.Value, nil
}

func (s *DirectClient) GetConfirmedBlockParsed(ctx context.Context, slot uint64) (GetConfirmBlockParsedResponse, error) {
	res := struct {
		GeneralResponse
		Result GetConfirmBlockParsedResponse `json:"result"`
	}{}
	err := s.request(ctx, "getConfirmedBlock", []interface{}{slot, "jsonParsed"}, &res)
	if err != nil {
		return GetConfirmBlockParsedResponse{}, err
	}
	return res.Result, nil
}

func (s *DirectClient) GetConfirmedTransactionParsed(ctx context.Context, txhash string) (GetConfirmedTransactionParsedResponse, error) {
	res := struct {
		GeneralResponse
		Result GetConfirmedTransactionParsedResponse `json:"result"`
	}{}
	err := s.request(ctx, "getConfirmedTransaction", []interface{}{txhash, "jsonParsed"}, &res)
	if err != nil {
		return GetConfirmedTransactionParsedResponse{}, err
	}
	return res.Result, nil
}

func (s *DirectClient) GetTokenAccountsByOwner(ctx context.Context, account string) (string, error) {
	res := struct {
		GeneralResponse
		Result struct {
			Context Context    `json:"context"`
			Value   []Accounts `json:"value"`
		} `json:"result"`
	}{}
	params := []interface{}{account,
		map[string]interface{}{"programId": TokenProgramID.ToBase58()},
		map[string]interface{}{
			"encoding": "jsonParsed",
		}}

	err := s.request(ctx, "getTokenAccountsByOwner", params, &res)
	if err != nil {
		return "", fmt.Errorf("No Token Account Found")
	}
	tokenAccounts := res.Result.Value
	if err != nil || len(tokenAccounts) == 0 {
		return "", fmt.Errorf("No Token Account Found")
	}
	return tokenAccounts[0].Pubkey, nil
}

func (s *DirectClient) GetTokenAccountByMint(ctx context.Context, account string, mint string) (string, error) {
	res := struct {
		GeneralResponse
		Result struct {
			Context Context    `json:"context"`
			Value   []Accounts `json:"value"`
		} `json:"result"`
	}{}
	params := []interface{}{account,
		map[string]interface{}{"mint": mint},
		map[string]interface{}{
			"encoding": "jsonParsed",
		}}

	err := s.request(ctx, "getTokenAccountsByOwner", params, &res)
	if err != nil {
		return "", fmt.Errorf("No Token Account Found")
	}
	tokenAccounts := res.Result.Value
	if err != nil || len(tokenAccounts) == 0 {
		return "", fmt.Errorf("No Token Account Found")
	}
	return tokenAccounts[0].Pubkey, nil
}