module github.com/imerkle/rosetta-solana-go

require (
	contrib.go.opencensus.io/exporter/stackdriver v0.13.14 // indirect
	github.com/coinbase/rosetta-sdk-go v0.8.3
	github.com/coinbase/rosetta-sdk-go/types v1.0.0
	github.com/fatih/color v1.16.0
	github.com/iancoleman/strcase v0.3.0
	github.com/mitchellh/copystructure v1.2.0
	github.com/mr-tron/base58 v1.2.0
	github.com/portto/solana-go-sdk v1.26.0
	github.com/spf13/cobra v1.8.0
	github.com/streamingfast/binary v0.0.0-20210928223119-44fc44e4a0b5
	github.com/streamingfast/logging v0.0.0-20211221170249-09a6ecb200a0 // indirect
	github.com/streamingfast/solana-go v0.5.1
	github.com/teris-io/shortid v0.0.0-20220617161101-71ec9f2aa569 // indirect
	github.com/test-go/testify v1.1.4
	github.com/tidwall/gjson v1.17.0 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/sync v0.6.0
	gotest.tools v2.2.0+incompatible
)

go 1.15

replace github.com/portto/solana-go-sdk => github.com/imerkle/solana-go-sdk v0.0.10

//replace github.com/portto/solana-go-sdk => ../solana-go-sdk
