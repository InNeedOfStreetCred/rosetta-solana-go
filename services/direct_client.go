package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	ss "github.com/portto/solana-go-sdk/client"
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

func NewClient(endpoint string) *DirectClient {
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
	Blockhash     string           `json:"blockhash"`
	FeeCalculator ss.FeeCalculator `json:"feeCalculator"`
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
	Authority     string        `json:"authority"`
	BlockHash     string        `json:"blockhash"`
	FeeCalculator FeeCalculator `json:"feeCalculator"`
}
type Parsed struct {
	Info Info `json:"info"`
}
type AccData struct {
	Parsed Parsed `json:"parsed"`
}

type FeeCalculator struct {
	LamportsPerSignature string `json:"lamportsPerSignature"`
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
