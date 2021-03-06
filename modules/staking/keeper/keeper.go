package keeper

import (
	"container/list"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/params"
	"github.com/netcloth/netcloth-chain/modules/staking/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const aminoCacheSize = 500

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Implements DelegationSet interface
var _ types.DelegationSet = Keeper{}

// keeper of the staking store
type Keeper struct {
	storeKey           sdk.StoreKey
	storeTKey          sdk.StoreKey
	cdc                *codec.Codec
	supplyKeeper       types.SupplyKeeper
	hooks              types.StakingHooks
	paramstore         params.Subspace
	validatorCache     map[string]cachedValidator
	validatorCacheList *list.List
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(cdc *codec.Codec, key, tkey sdk.StoreKey, supplyKeeper types.SupplyKeeper,
	paramstore params.Subspace) Keeper {

	// ensure bonded and not bonded module accounts are set
	if addr := supplyKeeper.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := supplyKeeper.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	return Keeper{
		storeKey:           key,
		storeTKey:          tkey,
		cdc:                cdc,
		supplyKeeper:       supplyKeeper,
		paramstore:         paramstore.WithKeyTable(ParamKeyTable()),
		hooks:              nil,
		validatorCache:     make(map[string]cachedValidator, aminoCacheSize),
		validatorCacheList: list.New(),
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("modules/%s", types.ModuleName))
}

// Set the validator hooks
func (k *Keeper) SetHooks(sh types.StakingHooks) *Keeper {
	if k.hooks != nil {
		panic("cannot set validator hooks twice")
	}
	k.hooks = sh
	return k
}

// Load the last total validator power.
func (k Keeper) GetLastTotalPower(ctx sdk.Context) (power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := store.Get(types.LastTotalPowerKey)
	if b == nil {
		return sdk.ZeroInt()
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &power)
	return
}

// Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshalBinaryLengthPrefixed(power)
	store.Set(types.LastTotalPowerKey, b)
}

func (k Keeper) EndBlock(ctx sdk.Context) {
	p := k.GetParams(ctx)

	if ctx.BlockTime().Unix() < p.NextExtendingTime {
		return
	}

	if p.MaxValidatorsExtending > p.MaxValidators {
		e := p.MaxValidatorsExtendingSpeed
		if p.MaxValidatorsExtending-p.MaxValidators < p.MaxValidatorsExtendingSpeed {
			e = p.MaxValidatorsExtending - p.MaxValidators
		}

		p.MaxValidators += e
	}

	p.NextExtendingTime += types.MaxValidatorsExtendingInterval
	k.SetParams(ctx, p)
}
