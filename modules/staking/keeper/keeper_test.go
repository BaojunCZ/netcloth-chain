package keeper_test

import (
	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/auth"
	"github.com/netcloth/netcloth-chain/modules/bank"
	"github.com/netcloth/netcloth-chain/modules/gov"
	"github.com/netcloth/netcloth-chain/modules/mint"
	"github.com/netcloth/netcloth-chain/modules/params"
	"github.com/netcloth/netcloth-chain/modules/staking"
	stakingtypes "github.com/netcloth/netcloth-chain/modules/staking/types"
	"github.com/netcloth/netcloth-chain/modules/supply"
	"github.com/netcloth/netcloth-chain/store"
	sdk "github.com/netcloth/netcloth-chain/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	"os"
	"testing"
	"time"

	distr "github.com/netcloth/netcloth-chain/modules/distribution"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	maccPerms = map[string][]string {
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            {supply.Burner},
	}
)
func moduleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	return modAccAddrs
}

func setupTest() (stakingKeeper staking.Keeper, ctx sdk.Context) {
	cdc := codec.New()

	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)

	keys := sdk.NewKVStoreKeys(params.StoreKey, auth.StoreKey, supply.StoreKey, staking.StoreKey)
	tkeys := sdk.NewTransientStoreKeys(params.TStoreKey, staking.TStoreKey)

	paramsKeeper := params.NewKeeper(cdc, keys[params.StoreKey], tkeys[params.TStoreKey], params.DefaultCodespace)

	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := paramsKeeper.Subspace(staking.DefaultParamspace)

	ms.MountStoreWithDB(keys[auth.StoreKey], sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeys[staking.TStoreKey], sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keys[staking.StoreKey], sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys[supply.StoreKey], sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keys[params.StoreKey], sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeys[params.TStoreKey], sdk.StoreTypeTransient, db)

	ms.LoadLatestVersion()

	accountKeeper := auth.NewAccountKeeper(cdc, keys[auth.StoreKey], authSubspace, auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, bankSubspace, bank.DefaultCodespace, moduleAccountAddrs())
	supplyKeeper := supply.NewKeeper(cdc, keys[supply.StoreKey], accountKeeper, bankKeeper, maccPerms)
	stakingKeeper = staking.NewKeeper(cdc, keys[staking.StoreKey], tkeys[staking.TStoreKey], supplyKeeper, stakingSubspace, staking.DefaultCodespace)
	ctx = sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, false, log.NewTMLogger(os.Stdout))

	return
}

func TestEndBlock(t *testing.T) {
	k, ctx := setupTest()

	p := staking.Params {
		MaxValidators                   : 100,
		MaxValidatorsExtending          : 130,
		MaxValidatorsExtendingSpeed     : 10,
		NextExtendingTime               : time.Now().Unix() + stakingtypes.MaxValidatorsExtendingInterval,
	}

	k.SetParams(ctx, p)

	require.Equal(t, 100, int(k.GetParams(ctx).MaxValidators))

	ctx = ctx.WithBlockTime(time.Now().Add(stakingtypes.MaxValidatorsExtendingInterval * 1e9 * 1))
	k.EndBlock(ctx)
	require.Equal(t, 110, int(k.GetParams(ctx).MaxValidators))

	p = k.GetParams(ctx)
	p.MaxValidatorsExtendingSpeed = 11
	k.SetParams(ctx, p)
	ctx = ctx.WithBlockTime(time.Now().Add(stakingtypes.MaxValidatorsExtendingInterval * 1e9 * 2))
	k.EndBlock(ctx)
	require.Equal(t, 121, int(k.GetParams(ctx).MaxValidators))

	ctx = ctx.WithBlockTime(time.Now().Add(stakingtypes.MaxValidatorsExtendingInterval * 1e9 * 3))
	k.EndBlock(ctx)
	require.Equal(t, 130, int(k.GetParams(ctx).MaxValidators))
}