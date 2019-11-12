package ipal

import (
	"github.com/netcloth/netcloth-chain/modules/ipal/keeper"
	"github.com/netcloth/netcloth-chain/modules/ipal/types"
)

const (
	ModuleName        = types.ModuleName
	StoreKey          = types.StoreKey
	RouterKey         = types.RouterKey
	QuerierRoute      = types.QuerierRoute
	DefaultCodespace  = types.DefaultCodespace
	DefaultParamspace = keeper.DefaultParamspace
)

var (
	NewKeeper              = keeper.NewKeeper
	NewQuerier             = keeper.NewQuerier
	RegisterCodec          = types.RegisterCodec
	NewServerNodeObject    = types.NewServiceNode
	NewMsgServiceNodeClaim = types.NewMsgServiceNodeClaim
	ErrEmptyInputs         = types.ErrEmptyInputs
	ModuleCdc              = types.ModuleCdc
	AttributeValueCategory = types.AttributeValueCategory
	NewEndpoint            = types.NewEndpoint
)

type (
	Keeper              = keeper.Keeper
	MsgServiceNodeClaim = types.MsgServiceNodeClaim
	Endpoint            = types.Endpoint
	Endpoints           = types.Endpoints
)
