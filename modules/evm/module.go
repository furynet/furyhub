package evm

import (
	"encoding/json"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	ethermint "github.com/evmos/ethermint/x/evm"
	"github.com/evmos/ethermint/x/evm/keeper"
	"github.com/evmos/ethermint/x/evm/types"

	furytypes "github.com/furynet/furyhub/types"
)

var (
	_ module.AppModule = AppModule{}
)

// ____________________________________________________________________________

// AppModule implements an application module for the evm module.
type AppModule struct {
	ethermint.AppModule
}

// NewAppModule creates a new AppModule object
func NewAppModule(k *keeper.Keeper, ak types.AccountKeeper) AppModule {
	return AppModule{
		AppModule: ethermint.NewAppModule(k, ak),
	}
}

// BeginBlock returns the begin block for the evm module.
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	ethChainID := furytypes.BuildEthChainID(ctx.ChainID())
	am.AppModule.BeginBlock(ctx.WithChainID(ethChainID), req)
}

// InitGenesis performs genesis initialization for the evm module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	ethChainID := furytypes.BuildEthChainID(ctx.ChainID())
	return am.AppModule.InitGenesis(ctx.WithChainID(ethChainID), cdc, data)
}
