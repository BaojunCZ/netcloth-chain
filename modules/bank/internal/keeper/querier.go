package keeper

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/bank/internal/types"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	// query balance path
	QueryBalance = "balances"
)

// NewQuerier returns a new sdk.Keeper instance.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryBalance:
			return queryBalance(ctx, req, k)

		default:
			return nil, sdk.ErrUnknownRequest("unknown bank query endpoint")
		}
	}
}

// queryBalance fetch an account's balance for the supplied height.
// Height and account address are passed as first and second path components respectively.
func queryBalance(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, sdk.Error) {
	var params types.QueryBalanceParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, k.GetCoins(ctx, params.Address))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
