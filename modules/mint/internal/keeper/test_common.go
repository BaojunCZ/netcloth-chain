// nolint:deadcode unused
package keeper

import (
	"os"
	"testing"
	"time"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/stretchr/testify/require"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/auth"
	"github.com/netcloth/netcloth-chain/modules/bank"
	"github.com/netcloth/netcloth-chain/modules/mint/internal/types"
	"github.com/netcloth/netcloth-chain/modules/params"
	"github.com/netcloth/netcloth-chain/modules/staking"
	"github.com/netcloth/netcloth-chain/modules/supply"
	"github.com/netcloth/netcloth-chain/store"
	sdk "github.com/netcloth/netcloth-chain/types"
)

type testInput struct {
	ctx        sdk.Context
	cdc        *codec.Codec
	mintKeeper Keeper
}

func newTestInput(t *testing.T) testInput {
	db := dbm.NewMemDB()

	keyAcc := sdk.NewKVStoreKey(auth.StoreKey)
	keySupply := sdk.NewKVStoreKey(supply.StoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyMint := sdk.NewKVStoreKey(types.StoreKey)

	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(keyAcc, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyStaking, sdk.StoreTypeTransient, nil)
	ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keySupply, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(keyMint, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(tkeyParams, sdk.StoreTypeTransient, db)
	err := ms.LoadLatestVersion()
	require.Nil(t, err)

	ctx := sdk.NewContext(ms, abci.Header{Time: time.Unix(0, 0)}, false, log.NewTMLogger(os.Stdout))

	feeCollectorAcc := supply.NewEmptyModuleAccount(auth.FeeCollectorName)
	notBondedPool := supply.NewEmptyModuleAccount(staking.NotBondedPoolName, supply.Burner, supply.Staking)
	bondPool := supply.NewEmptyModuleAccount(staking.BondedPoolName, supply.Burner, supply.Staking)
	minterAcc := supply.NewEmptyModuleAccount(types.ModuleName, supply.Minter)

	blacklistedAddrs := make(map[string]bool)
	blacklistedAddrs[feeCollectorAcc.String()] = true
	blacklistedAddrs[notBondedPool.String()] = true
	blacklistedAddrs[bondPool.String()] = true
	blacklistedAddrs[minterAcc.String()] = true

	paramsKeeper := params.NewKeeper(types.ModuleCdc, keyParams, tkeyParams)
	accountKeeper := auth.NewAccountKeeper(types.ModuleCdc, keyAcc, paramsKeeper.Subspace(auth.DefaultParamspace), auth.ProtoBaseAccount)
	bankKeeper := bank.NewBaseKeeper(accountKeeper, paramsKeeper.Subspace(bank.DefaultParamspace), blacklistedAddrs)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		types.ModuleName:          []string{supply.Minter},
		staking.NotBondedPoolName: []string{supply.Burner, supply.Staking},
		staking.BondedPoolName:    []string{supply.Burner, supply.Staking},
	}
	supplyKeeper := supply.NewKeeper(types.ModuleCdc, keySupply, accountKeeper, bankKeeper, maccPerms)
	supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{}))

	stakingKeeper := staking.NewKeeper(
		types.ModuleCdc, keyStaking, tkeyStaking, supplyKeeper, paramsKeeper.Subspace(staking.DefaultParamspace))
	mintKeeper := NewKeeper(types.ModuleCdc, keyMint, paramsKeeper.Subspace(types.DefaultParamspace), &stakingKeeper, supplyKeeper, auth.FeeCollectorName)

	// set module accounts
	supplyKeeper.SetModuleAccount(ctx, feeCollectorAcc)
	supplyKeeper.SetModuleAccount(ctx, minterAcc)
	supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	supplyKeeper.SetModuleAccount(ctx, bondPool)

	mintKeeper.SetParams(ctx, types.DefaultParams())

	return testInput{ctx, types.ModuleCdc, mintKeeper}
}
