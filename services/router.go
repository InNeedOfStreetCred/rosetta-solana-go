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
	"net/http"

	"github.com/imerkle/rosetta-solana-go/configuration"
	solanago "github.com/imerkle/rosetta-solana-go/solana"

	"github.com/coinbase/rosetta-sdk-go/asserter"
	"github.com/coinbase/rosetta-sdk-go/server"
)

// NewBlockchainRouter creates a Mux http.Handler from a collection
// of server controllers.
func NewBlockchainRouter(
	config *configuration.Configuration,
	client *solanago.Client,
	asserter *asserter.Asserter,
) http.Handler {
	networkAPIService := NewNetworkAPIService(config, client)
	networkAPIController := server.NewNetworkAPIController(
		networkAPIService,
		asserter,
	)

	accountAPIService := NewAccountAPIService(config, client)
	accountAPIController := server.NewAccountAPIController(
		accountAPIService,
		asserter,
	)

	blockAPIService := NewBlockAPIService(config, client)
	blockAPIController := server.NewBlockAPIController(
		blockAPIService,
		asserter,
	)

	constructionAPIService := NewConstructionAPIService(config, client)
	constructionAPIController := server.NewConstructionAPIController(
		constructionAPIService,
		asserter,
	)

	callAPIService := NewCallAPIService(config, client)
	callAPIController := server.NewCallAPIController(
		callAPIService,
		asserter,
	)

	return server.NewRouter(
		networkAPIController,
		accountAPIController,
		blockAPIController,
		constructionAPIController,
		callAPIController,
	)
}