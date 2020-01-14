package params

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"

	"github.com/netcloth/netcloth-chain/codec"
	"github.com/netcloth/netcloth-chain/modules/params/subspace"
	"github.com/netcloth/netcloth-chain/modules/params/types"
	"github.com/netcloth/netcloth-chain/store"
	sdk "github.com/netcloth/netcloth-chain/types"
)

func validateNoOp(_ interface{}) error { return nil }

type testInput struct {
	ctx    sdk.Context
	cdc    *codec.Codec
	keeper Keeper
}

var (
	_ subspace.ParamSet = (*testParams)(nil)

	keyMaxValidators = "MaxValidators"
	keySlashingRate  = "SlashingRate"
	testSubspace     = "TestSubspace"
)

type testParamsSlashingRate struct {
	DoubleSign uint16 `json:"double_sign,omitempty" yaml:"double_sign,omitempty"`
	Downtime   uint16 `json:"downtime,omitempty" yaml:"downtime,omitempty"`
}

type testParams struct {
	MaxValidators uint16                 `json:"max_validators" yaml:"max_validators"` // maximum number of validators (max uint16 = 65535)
	SlashingRate  testParamsSlashingRate `json:"slashing_rate" yaml:"slashing_rate"`
}

func (tp *testParams) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		NewParamSetPair([]byte(keyMaxValidators), &tp.MaxValidators, validateNoOp),
		NewParamSetPair([]byte(keySlashingRate), &tp.SlashingRate, validateNoOp),
	}
}

func testProposal(changes ...ParamChange) ParameterChangeProposal {
	return NewParameterChangeProposal(
		"Test",
		"description",
		changes,
	)
}

func newTestInput(t *testing.T) testInput {
	cdc := codec.New()
	types.RegisterCodec(cdc)

	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	keyParams := sdk.NewKVStoreKey("params")
	tKeyParams := sdk.NewTransientStoreKey("transient_params")

	cms.MountStoreWithDB(keyParams, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tKeyParams, sdk.StoreTypeTransient, db)

	err := cms.LoadLatestVersion()
	require.Nil(t, err)

	keeper := NewKeeper(cdc, keyParams, tKeyParams)
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())

	return testInput{ctx, cdc, keeper}
}

func TestProposalHandlerPassed(t *testing.T) {
	input := newTestInput(t)
	ss := input.keeper.Subspace(testSubspace).WithKeyTable(
		NewKeyTable().RegisterParamSet(&testParams{}),
	)

	tp := testProposal(NewParamChange(testSubspace, keyMaxValidators, "1"))
	hdlr := NewParamChangeProposalHandler(input.keeper)
	require.NoError(t, hdlr(input.ctx, tp))

	var param uint16
	ss.Get(input.ctx, []byte(keyMaxValidators), &param)
	require.Equal(t, param, uint16(1))
}

func TestProposalHandlerFailed(t *testing.T) {
	input := newTestInput(t)
	ss := input.keeper.Subspace(testSubspace).WithKeyTable(
		NewKeyTable().RegisterParamSet(&testParams{}),
	)

	tp := testProposal(NewParamChange(testSubspace, keyMaxValidators, "invalidType"))
	hdlr := NewParamChangeProposalHandler(input.keeper)
	require.Error(t, hdlr(input.ctx, tp))

	require.False(t, ss.Has(input.ctx, []byte(keyMaxValidators)))
}

func TestProposalHandlerUpdateOmitempty(t *testing.T) {
	input := newTestInput(t)
	ss := input.keeper.Subspace(testSubspace).WithKeyTable(
		NewKeyTable().RegisterParamSet(&testParams{}),
	)

	hdlr := NewParamChangeProposalHandler(input.keeper)
	var param testParamsSlashingRate

	tp := testProposal(NewParamChange(testSubspace, keySlashingRate, `{"downtime": 7}`))
	require.NoError(t, hdlr(input.ctx, tp))

	ss.Get(input.ctx, []byte(keySlashingRate), &param)
	require.Equal(t, testParamsSlashingRate{0, 7}, param)

	tp = testProposal(NewParamChange(testSubspace, keySlashingRate, `{"double_sign": 10}`))
	require.NoError(t, hdlr(input.ctx, tp))

	ss.Get(input.ctx, []byte(keySlashingRate), &param)
	require.Equal(t, testParamsSlashingRate{10, 7}, param)
}