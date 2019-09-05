package types

import (
	"github.com/NetCloth/netcloth-chain/codec"
)

// Register concrete types on codec codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateValidator{}, "nch/MsgCreateValidator", nil)
	cdc.RegisterConcrete(MsgEditValidator{}, "nch/MsgEditValidator", nil)
	cdc.RegisterConcrete(MsgDelegate{}, "nch/MsgDelegate", nil)
	cdc.RegisterConcrete(MsgUndelegate{}, "nch/MsgUndelegate", nil)
	cdc.RegisterConcrete(MsgBeginRedelegate{}, "nch/MsgBeginRedelegate", nil)
}

// generic sealed codec to be used throughout this module
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
