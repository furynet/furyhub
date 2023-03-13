package guardian_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/furynet/furyhub/modules/guardian"
	"github.com/furynet/furyhub/modules/guardian/keeper"
	"github.com/furynet/furyhub/modules/guardian/types"
	"github.com/furynet/furyhub/simapp"
)

type TestSuite struct {
	suite.Suite

	cdc    codec.Codec
	ctx    sdk.Context
	keeper keeper.Keeper
}

func (suite *TestSuite) SetupTest() {
	app := simapp.Setup(suite.T(), false)

	suite.cdc = codec.NewAminoCodec(app.LegacyAmino())
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{})
	suite.keeper = app.GuardianKeeper
}

func TestGenesisSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestExportGenesis() {
	exportedGenesis := guardian.ExportGenesis(suite.ctx, suite.keeper)
	defaultGenesis := types.DefaultGenesisState()
	suite.Equal(exportedGenesis, defaultGenesis)
}
