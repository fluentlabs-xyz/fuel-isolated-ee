package graphql_entrypoints

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/fluentlabs-xyz/fuel-ee/graphql_object"
	"github.com/fluentlabs-xyz/fuel-ee/graphql_scalars"
	"github.com/fluentlabs-xyz/fuel-ee/types"
	"github.com/graphql-go/graphql"
	"log"
)

type DryRunEntry struct {
	SchemaFields graphql_object.SchemaFields
}

type DryRunEntryStruct struct {
}

// const encodedTransactionsArgName = "encodedTransactions"
const encodedTransactionsArgName = "txs"
const utxoValidationArgName = "utxoValidation"
const gasPriceArgName = "gasPrice"

func MakeDryRunEntry(ethClient *ethclient.Client, dryRunTransactionStatusType *graphql_object.DryRunTransactionExecutionStatusType) (*DryRunEntry, error) {
	objectConfig := graphql.ObjectConfig{Name: "DryRunEntry", Fields: graphql.Fields{
		"dryRun": &graphql.Field{
			Type: graphql.NewList(dryRunTransactionStatusType.SchemaFields.Object),
			Args: graphql.FieldConfigArgument{
				encodedTransactionsArgName: &graphql.ArgumentConfig{
					Type:         graphql.NewList(graphql_scalars.HexStringType),
					DefaultValue: []graphql_scalars.HexString{},
				},
				utxoValidationArgName: &graphql.ArgumentConfig{
					Type:         graphql.Boolean,
					DefaultValue: false,
				},
				gasPriceArgName: &graphql.ArgumentConfig{
					Type:         graphql_scalars.U64Type,
					DefaultValue: graphql_scalars.NewU64(0),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				encodedTransactions := p.Args[encodedTransactionsArgName]
				//utxoValidation := p.Args[utxoValidationArgName]
				//gasPrice := p.Args[gasPriceArgName]
				encodedTransactionsList, ok := encodedTransactions.([]interface{})
				if !ok {
					return nil, errors.New("encoded transactions must be a list")
				}
				for _, encodedTransaction := range encodedTransactionsList {
					transactionHexString, ok := encodedTransaction.(*graphql_scalars.HexString)
					if !ok {
						return nil, errors.New("each encoded transaction must be a hex string")
					}
					log.Printf("transactionHexString: %s", transactionHexString)

					fvmDryRunSigBytes := make([]byte, 4)
					binary.BigEndian.PutUint32(fvmDryRunSigBytes, types.FvmDryRunSig)
					// send tx to reth node for emulation/estimation process (to get status, receipts, gas spent)
					data := append(fvmDryRunSigBytes, transactionHexString.Value()...)
					from := common.HexToAddress(types.FuelRelayerAccountAddress)
					to := common.HexToAddress(types.EthFuelVMPrecompileAddress)
					callMsg := ethereum.CallMsg{
						From: from,
						To:   &to,
						Data: data,
					}
					estimatedGas, err := ethClient.EstimateGas(context.Background(), callMsg)
					if err != nil {
						return nil, errors.New(fmt.Sprintf("dryRun: failed to estimate gas, error: %s", err))
					}
					log.Printf("estimatedGas: %d", estimatedGas)
					callMsg.Gas = estimatedGas
					callRes, err := ethClient.CallContract(context.Background(), callMsg, nil)
					if err != nil {
						return nil, errors.New(fmt.Sprintf("dryRun: failed to call contract, error: %s", err))
					}
					log.Printf("callRes: %s", callRes)
				}
				return []graphql_object.DryRunTransactionExecutionStatusStruct{
					{
						Id:       "0xb4f5b359704eda15f8ec6c15004b6816b9df4f730baaa50d0a2fb34a99108bee",
						Status:   &graphql_object.DryRunSuccessStatusStruct{},
						Receipts: []graphql_object.ReceiptStruct{},
					},
				}, nil
			},
		},
	}}
	object := graphql.NewObject(objectConfig)
	schemaConfig := graphql.SchemaConfig{
		Query:    object,
		Mutation: object,
	}
	schema, err := graphql.NewSchema(schemaConfig)

	return &DryRunEntry{
		SchemaFields: graphql_object.SchemaFields{
			Schema:       &schema,
			ObjectConfig: &objectConfig,
			Object:       object,
			SchemaConfig: &schemaConfig,
		},
	}, err
}
