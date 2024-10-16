package graphql_object

import (
	"github.com/fluentlabs-xyz/fuel-ee/src/config"
	"github.com/fluentlabs-xyz/fuel-ee/src/graphql_scalars"
	"github.com/graphql-go/graphql"
)

type ConsensusParametersType struct {
	SchemaFields SchemaFields
}

//	pub struct ConsensusParameters {
//	   pub version: ConsensusParametersVersion,
//	   pub tx_params: TxParameters,
//	   pub predicate_params: PredicateParameters,
//	   pub script_params: ScriptParameters,
//	   pub contract_params: ContractParameters,
//	   pub fee_params: FeeParameters,
//	   pub base_asset_id: AssetId,
//	   pub block_gas_limit: U64,
//	   pub chain_id: U64,
//	   pub gas_costs: GasCosts,
//	   pub privileged_address: Address,
//	}
type ConsensusParametersStruct struct {
}

func ConsensusParameters(
	config *config.Config,
	consensusParametersVersionType *ConsensusParametersVersionType,
	txParametersType *TxParametersType,
	predicateParametersType *PredicateParametersType,
	scriptParametersType *ScriptParametersType,
	contractParametersType *ContractParametersType,
	feeParametersType *FeeParametersType,
	gasCostsType *GasCostsType,
) (*ConsensusParametersType, error) {
	objectConfig := graphql.ObjectConfig{Name: "ConsensusParameters", Fields: graphql.Fields{
		"version": &graphql.Field{
			Type: consensusParametersVersionType.SchemaFields.Type,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return 1, nil
			},
		},
		"chainId": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return config.Blockchain.ChainId, nil
			},
		},
		"txParams": &graphql.Field{
			Type: txParametersType.SchemaFields.Object,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return &TxParametersVersionStruct{}, nil
			},
		},
		"predicateParams": &graphql.Field{
			Type: predicateParametersType.SchemaFields.Object,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return &PredicateParametersStruct{}, nil
			},
		},
		"scriptParams": &graphql.Field{
			Type: scriptParametersType.SchemaFields.Object,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return &ScriptParametersStruct{
					MaxScriptLength:     102400,
					MaxScriptDataLength: 102400,
				}, nil
			},
		},
		"contractParams": &graphql.Field{
			Type: contractParametersType.SchemaFields.Object,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return &ContractParametersStruct{
					ContractMaxSize: 102400,
					MaxStorageSlots: 1760,
				}, nil
			},
		},
		"feeParams": &graphql.Field{
			Type: feeParametersType.SchemaFields.Object,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return &FeeParametersStruct{}, nil
			},
		},
		"gasCosts": &graphql.Field{
			Type: gasCostsType.SchemaFields.Object,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return &GasCostsStruct{
					aloc: &LightOperationStruct{
						Base:        0,
						UnitsPerGas: 0,
					},
				}, nil
			},
		},
		"baseAssetId": &graphql.Field{
			Type: graphql_scalars.Bytes32Type,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return graphql_scalars.NewBytes32TryFromStringOrPanic(config.Blockchain.FuelBaseAssetId), nil
			},
		},
	}}
	object := graphql.NewObject(objectConfig)
	schemaConfig := graphql.SchemaConfig{Query: object}
	schema, err := graphql.NewSchema(schemaConfig)

	return &ConsensusParametersType{
		SchemaFields: SchemaFields{
			Schema:       &schema,
			ObjectConfig: &objectConfig,
			Object:       object,
			SchemaConfig: &schemaConfig,
		},
	}, err
}
