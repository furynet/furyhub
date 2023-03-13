package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/furynet/furyhub/modules/mint/types"
	"github.com/furynet/furyhub/simapp"
)

type KeeperTestSuite struct {
	suite.Suite

	cdc *codec.LegacyAmino
	ctx sdk.Context
	app *simapp.SimApp
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(suite.T(), false)

	suite.cdc = app.LegacyAmino()
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	suite.app = app

	app.MintKeeper.SetParamSet(suite.ctx, types.DefaultParams())
	app.MintKeeper.SetMinter(suite.ctx, types.DefaultMinter())
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestSetGetMinter() {
	minter := types.NewMinter(time.Now().UTC(), sdk.NewInt(100000))
	suite.app.MintKeeper.SetMinter(suite.ctx, minter)
	expMinter := suite.app.MintKeeper.GetMinter(suite.ctx)

	require.Equal(suite.T(), minter, expMinter)
}

func (suite *KeeperTestSuite) TestSetGetParamSet() {
	suite.app.MintKeeper.SetParamSet(suite.ctx, types.DefaultParams())
	expParamSet := suite.app.MintKeeper.GetParamSet(suite.ctx)

	require.Equal(suite.T(), types.DefaultParams(), expParamSet)
}

func (suite *KeeperTestSuite) TestMintCoins() {

	mintCoins := sdk.NewCoins(sdk.NewCoin("fury", sdk.NewInt(1000)))
	err := suite.app.MintKeeper.MintCoins(suite.ctx, mintCoins)
	require.NoError(suite.T(), err)

	acc := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	coins := suite.app.BankKeeper.GetAllBalances(suite.ctx, acc.GetAddress())
	require.Equal(suite.T(), coins, mintCoins)
}

func (suite *KeeperTestSuite) TestAddCollectedFees() {
	mintCoins := sdk.NewCoins(sdk.NewCoin("fury", sdk.NewInt(1000)))

	feeCollector := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, "fee_collector")
	feeCollectorBalance := suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector.GetAddress())

	err := suite.app.MintKeeper.MintCoins(suite.ctx, mintCoins)
	require.NoError(suite.T(), err)

	acc := suite.app.AccountKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	coins := suite.app.BankKeeper.GetAllBalances(suite.ctx, acc.GetAddress())
	require.Equal(suite.T(), coins, mintCoins)

	err = suite.app.MintKeeper.AddCollectedFees(suite.ctx, mintCoins)
	require.NoError(suite.T(), err)

	acc = suite.app.AccountKeeper.GetModuleAccount(suite.ctx, types.ModuleName)
	coins = suite.app.BankKeeper.GetAllBalances(suite.ctx, acc.GetAddress())
	require.True(suite.T(), coins.Empty())

	feeCollectorTotalBalance := suite.app.BankKeeper.GetAllBalances(suite.ctx, feeCollector.GetAddress())
	expectedCollectedFees := feeCollectorTotalBalance.Sub(feeCollectorBalance...)
	require.Equal(suite.T(), expectedCollectedFees, mintCoins)

}
