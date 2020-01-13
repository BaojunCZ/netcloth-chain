package cipal

import (
	"github.com/netcloth/netcloth-chain/modules/cipal/keeper"
	"github.com/netcloth/netcloth-chain/modules/cipal/types"
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
	RegisterCodec                = types.RegisterCodec
	NewIPALObject                = types.NewCIPALObject
	NewQuerier                   = keeper.NewQuerier
	NewParam                     = types.NewParam
	NewUserRequest               = types.NewUserRequest
	NewMsgIPAL                   = types.NewMsgCIPAL
	NewKeeper                    = keeper.NewKeeper
	ErrEmptyInputs               = types.ErrEmptyInputs
	ErrStringTooLong             = types.ErrStringTooLong
	ErrInvalidSignature          = types.ErrInvalidSignature
	ErrCIPALUserRequestExpired   = types.ErrCIPALUserRequestExpired
	ErrCIPALUserRequestSigVerify = types.ErrCIPALUserRequestSigVerify
	ModuleCdc                    = types.ModuleCdc
	AttributeValueCategory       = types.AttributeValueCategory
)

type (
	Keeper          = keeper.Keeper
	MsgIPAL         = types.MsgCIPAL
	IPALUserRequest = types.UserRequest
	Param           = types.Param
	CIPALObject     = types.CIPALObject
)
