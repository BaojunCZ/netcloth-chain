package types

import (
	sdk "github.com/netcloth/netcloth-chain/types"
)

const (
	DefaultCodespace sdk.CodespaceType = ModuleName

	CodeEmptyInputs               sdk.CodeType = 110
	CodeStringTooLong             sdk.CodeType = 111
	CodeInvalidIPALUserRequestSig sdk.CodeType = 112
	CodeCIPALUserRequestExpired   sdk.CodeType = 113
	CodeCIPALUserRequestSigVerify sdk.CodeType = 114
)

func ErrEmptyInputs(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeEmptyInputs, msg)
}

func ErrStringTooLong(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeStringTooLong, msg)
}

func ErrInvalidSignature(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidIPALUserRequestSig, msg)
}

func ErrCIPALUserRequestExpired(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeCIPALUserRequestExpired, msg)
}

func ErrCIPALUserRequestSigVerify(msg string) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeCIPALUserRequestSigVerify, msg)
}
