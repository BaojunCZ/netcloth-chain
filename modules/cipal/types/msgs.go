package types

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/netcloth/netcloth-chain/modules/auth"
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	maxUserAddressLength   = 256
	maxServerAddressLength = 256
)

var (
	_ sdk.Msg = MsgCIPAL{}
)

type ServiceInfo struct {
	Type    uint64 `json:"type" yaml:"type"`
	Address string `json:"address" yaml:"address"`
}

type Param struct {
	UserAddress string      `json:"user_address" yaml:"user_address"`
	ServiceInfo ServiceInfo `json:"service_info" yaml:"service_info"`
	Expiration  time.Time   `json:"expiration"`
}

type UserRequest struct {
	Params Param             `json:"params" yaml:"params"`
	Sig    auth.StdSignature `json:"signature" yaml:"signature`
}

type MsgCIPAL struct {
	From        sdk.AccAddress `json:"from" yaml:"from`
	UserRequest UserRequest    `json:"user_request" yaml:"user_request"`
}

func (i ServiceInfo) Validate() sdk.Error {
	if i.Address == "" {
		return ErrEmptyInputs("server address empty")
	}

	if len(i.Address) > maxServerAddressLength {
		return ErrStringTooLong("server address too long")
	}

	return nil
}

func (i ServiceInfo) String() string {
	return fmt.Sprintf(`ServiceInfo{Type:%s,Address:%s`, i.Type, i.Address)
}

func (p Param) GetSignBytes() []byte {
	b, err := json.Marshal(p)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (p Param) Validate() sdk.Error {
	if p.UserAddress == "" {
		return ErrEmptyInputs("user address empty")
	}

	if len(p.UserAddress) > maxUserAddressLength {
		return ErrStringTooLong("user address too long")
	}

	return p.ServiceInfo.Validate()
}

func NewParam(userAddress string, serviceAddress string, serviceType uint64, expiration time.Time) Param {
	return Param{
		UserAddress: userAddress,
		ServiceInfo: ServiceInfo{Type: serviceType, Address: serviceAddress},
		Expiration:  expiration,
	}
}

func NewUserRequest(userAddress string, serviceAddress string, serviceType uint64, expiration time.Time, sig auth.StdSignature) UserRequest {
	return UserRequest{
		Params: NewParam(userAddress, serviceAddress, serviceType, expiration),
		Sig:    sig,
	}
}

func NewMsgCIPAL(from sdk.AccAddress, userAddress string, serviceAddress string, serviceType uint64, expiration time.Time, sig auth.StdSignature) MsgCIPAL {
	return MsgCIPAL{
		from,
		NewUserRequest(userAddress, serviceAddress, serviceType, expiration, sig),
	}
}

func (msg MsgCIPAL) Route() string { return RouterKey }

func (msg MsgCIPAL) Type() string { return "cipal" }

func (msg MsgCIPAL) ValidateBasic() sdk.Error {
	if msg.From.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}

	err := msg.UserRequest.Params.Validate()
	if err != nil {
		return err
	}

	pubKey := msg.UserRequest.Sig.PubKey
	signBytes := msg.UserRequest.Params.GetSignBytes()
	if !pubKey.VerifyBytes(signBytes, msg.UserRequest.Sig.Signature) {
		return ErrInvalidSignature("user request signature invalid")
	}

	return nil
}

func (msg MsgCIPAL) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCIPAL) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}
