package app

import (
	sdkmath "cosmossdk.io/math"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	sdkupgrade "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	ica "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts"
	icacontrollertypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/controller/types"
	icahosttypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v5/modules/apps/27-interchain-accounts/types"
	ibchost "github.com/cosmos/ibc-go/v5/modules/core/24-host"

	coinswaptypes "github.com/irisnet/irismod/modules/coinswap/types"
	farmtypes "github.com/irisnet/irismod/modules/farm/types"
	"github.com/irisnet/irismod/modules/htlc"
	htlctypes "github.com/irisnet/irismod/modules/htlc/types"
	mttypes "github.com/irisnet/irismod/modules/mt/types"
	nftmodule "github.com/irisnet/irismod/modules/nft/module"
	nfttypes "github.com/irisnet/irismod/modules/nft/types"
	"github.com/irisnet/irismod/modules/oracle"
	oracletypes "github.com/irisnet/irismod/modules/oracle/types"
	"github.com/irisnet/irismod/modules/random"
	randomtypes "github.com/irisnet/irismod/modules/random/types"
	"github.com/irisnet/irismod/modules/record"
	recordtypes "github.com/irisnet/irismod/modules/record/types"
	"github.com/irisnet/irismod/modules/service"
	servicetypes "github.com/irisnet/irismod/modules/service/types"
	"github.com/irisnet/irismod/modules/token"
	tokentypes "github.com/irisnet/irismod/modules/token/types"

	tibcmttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/mt_transfer/types"
	tibcnfttypes "github.com/bianjieai/tibc-go/modules/tibc/apps/nft_transfer/types"
	tibcclienttypes "github.com/bianjieai/tibc-go/modules/tibc/core/02-client/types"
	tibchost "github.com/bianjieai/tibc-go/modules/tibc/core/24-host"

	migratehtlc "github.com/furynet/furyhub/migrate/htlc"
	migrateservice "github.com/furynet/furyhub/migrate/service"
	migratetibc "github.com/furynet/furyhub/migrate/tibc"
	"github.com/furynet/furyhub/modules/guardian"
	guardiantypes "github.com/furynet/furyhub/modules/guardian/types"
	"github.com/furynet/furyhub/modules/mint"
	minttypes "github.com/furynet/furyhub/modules/mint/types"
)

// RegisterUpgradePlan register a handler of upgrade plan
func (app *GridApp) RegisterUpgradePlan(cfg module.Configurator) {
	app.RegisterUpgradeHandler(
		"v1.1", &store.StoreUpgrades{},
		func(ctx sdk.Context, plan sdkupgrade.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// migrate htlc
			if err := migratehtlc.Migrate(ctx, app.appCodec, app.HTLCKeeper, app.BankKeeper, app.keys[htlctypes.StoreKey]); err != nil {
				panic(err)
			}
			// migrate service
			if err := migrateservice.Migrate(ctx, app.ServiceKeeper, app.BankKeeper); err != nil {
				panic(err)
			}

			return fromVM, nil
		},
	)
	app.RegisterUpgradeHandler(
		"v1.2", &store.StoreUpgrades{
			Added: []string{farmtypes.StoreKey, feegrant.StoreKey, tibchost.StoreKey, tibcnfttypes.StoreKey},
		},
		func(ctx sdk.Context, plan sdkupgrade.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// init farm params
			amount := sdkmath.NewIntWithDecimal(1000, int(nativeToken.Scale))
			farmtypes.SetDefaultGenesisState(farmtypes.GenesisState{
				Params: farmtypes.Params{
					PoolCreationFee:     sdk.NewCoin(nativeToken.MinUnit, amount),
					MaxRewardCategories: 2,
				}},
			)
			tibcclienttypes.SetDefaultGenesisState(tibcclienttypes.GenesisState{
				NativeChainName: "gridiron-mainnet",
			})

			if err := migratetibc.CreateClient(ctx,
				app.appCodec,
				"v1.2",
				app.TIBCKeeper.ClientKeeper,
			); err != nil {
				return nil, err
			}
			fromVM[authtypes.ModuleName] = 1
			fromVM[banktypes.ModuleName] = 1
			fromVM[stakingtypes.ModuleName] = 1
			fromVM[govtypes.ModuleName] = 1
			fromVM[distrtypes.ModuleName] = 1
			fromVM[slashingtypes.ModuleName] = 1
			fromVM[coinswaptypes.ModuleName] = 1
			fromVM[ibchost.ModuleName] = 1
			fromVM[capabilitytypes.ModuleName] = capability.AppModule{}.ConsensusVersion()
			fromVM[genutiltypes.ModuleName] = genutil.AppModule{}.ConsensusVersion()
			fromVM[minttypes.ModuleName] = mint.AppModule{}.ConsensusVersion()
			fromVM[paramstypes.ModuleName] = params.AppModule{}.ConsensusVersion()
			fromVM[crisistypes.ModuleName] = crisis.AppModule{}.ConsensusVersion()
			fromVM[upgradetypes.ModuleName] = crisis.AppModule{}.ConsensusVersion()
			fromVM[evidencetypes.ModuleName] = evidence.AppModule{}.ConsensusVersion()
			fromVM[feegrant.ModuleName] = feegrantmodule.AppModule{}.ConsensusVersion()
			fromVM[guardiantypes.ModuleName] = guardian.AppModule{}.ConsensusVersion()
			fromVM[tokentypes.ModuleName] = token.AppModule{}.ConsensusVersion()
			fromVM[recordtypes.ModuleName] = record.AppModule{}.ConsensusVersion()
			fromVM[nfttypes.ModuleName] = nftmodule.AppModule{}.ConsensusVersion()
			fromVM[htlctypes.ModuleName] = htlc.AppModule{}.ConsensusVersion()
			fromVM[servicetypes.ModuleName] = service.AppModule{}.ConsensusVersion()
			fromVM[oracletypes.ModuleName] = oracle.AppModule{}.ConsensusVersion()
			fromVM[randomtypes.ModuleName] = random.AppModule{}.ConsensusVersion()
			return app.mm.RunMigrations(ctx, cfg, fromVM)
		},
	)

	app.RegisterUpgradeHandler("v1.3",
		&store.StoreUpgrades{
			Added: []string{tibcmttypes.StoreKey, mttypes.StoreKey},
		},
		func(ctx sdk.Context, plan sdkupgrade.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			if err := migratetibc.CreateClient(ctx,
				app.appCodec,
				"v1.3",
				app.TIBCKeeper.ClientKeeper,
			); err != nil {
				return nil, err
			}
			return app.mm.RunMigrations(ctx, cfg, fromVM)
		},
	)

	app.RegisterUpgradeHandler("v1.4",
		&store.StoreUpgrades{
			Added: []string{authzkeeper.StoreKey},
		},
		func(ctx sdk.Context, plan sdkupgrade.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			// version upgrade:
			//	nft :    1 -> 2
			// 	auth:    2 -> 3
			// 	bank:    2 -> 3
			//	coinswap 3 -> 4
			// 	feegrant 1 -> 2
			// 	gov      2 -> 3
			// 	staking  2 -> 3
			// 	upgrade  2 -> 3

			// added module:
			//  authz

			// ibc application:
			//  27-interchain-accounts
			icaModule := app.mm.Modules[icatypes.ModuleName].(ica.AppModule)
			fromVM[icatypes.ModuleName] = icaModule.ConsensusVersion()
			// create ICS27 Controller submodule params
			controllerParams := icacontrollertypes.Params{}
			// create ICS27 Host submodule params
			hostParams := icahosttypes.Params{
				HostEnabled: true,
				AllowMessages: []string{
					authzMsgExec,
					authzMsgGrant,
					authzMsgRevoke,
					bankMsgSend,
					bankMsgMultiSend,
					distrMsgSetWithdrawAddr,
					distrMsgWithdrawValidatorCommission,
					distrMsgFundCommunityPool,
					distrMsgWithdrawDelegatorReward,
					feegrantMsgGrantAllowance,
					feegrantMsgRevokeAllowance,
					legacyGovMsgVoteWeighted,
					legacyGovMsgSubmitProposal,
					legacyGovMsgDeposit,
					legacyGovMsgVote,
					govMsgVoteWeighted,
					govMsgSubmitProposal,
					govMsgDeposit,
					govMsgVote,
					stakingMsgEditValidator,
					stakingMsgDelegate,
					stakingMsgUndelegate,
					stakingMsgBeginRedelegate,
					stakingMsgCreateValidator,
					vestingMsgCreateVestingAccount,
					ibcMsgTransfer,

					nftMsgIssueDenom,
					nftMsgTransferDenom,
					nftMsgMintNFT,
					nftMsgEditNFT,
					nftMsgTransferNFT,
					nftMsgBurnNFT,

					mtMsgIssueDenom,
					mtMsgTransferDenom,
					mtMsgMintMT,
					mtMsgEditMT,
					mtMsgTransferMT,
					mtMsgBurnMT,
				},
			}

			ctx.Logger().Info("start to init interchainaccount module...")
			// initialize ICS27 module
			icaModule.InitModule(ctx, controllerParams, hostParams)
			ctx.Logger().Info("start to run module migrations...")
			return app.mm.RunMigrations(ctx, cfg, fromVM)
		},
	)
}
