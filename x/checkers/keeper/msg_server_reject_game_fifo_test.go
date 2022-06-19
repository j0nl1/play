package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/j0nl1/play/x/checkers/types"
	"github.com/stretchr/testify/require"
)

func TestRejectSecondGameHasSavedFifo(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForRejectGame(t)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Red:     carol,
		Black:   alice,
	})
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator: carol,
		IdValue: "1",
	})
	nextGame, found := keeper.GetNextGame(sdk.UnwrapSDKContext(context))
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		IdValue:  3,
		FifoHead: "2",
		FifoTail: "2",
	}, nextGame)
	game2, found2 := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "2")
	require.True(t, found2)
	require.EqualValues(t, types.StoredGame{
		Creator:   bob,
		Index:     "2",
		Game:      "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       carol,
		Black:     alice,
		MoveCount: uint64(0),
		BeforeId:  "-1",
		AfterId:   "-1",
	}, game2)
}

func TestRejectMiddleGameHasSavedFifo(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForRejectGame(t)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Red:     carol,
		Black:   alice,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Red:     alice,
		Black:   bob,
	})
	msgServer.RejectGame(context, &types.MsgRejectGame{
		Creator: carol,
		IdValue: "2",
	})
	nextGame, found := keeper.GetNextGame(sdk.UnwrapSDKContext(context))
	require.True(t, found)
	require.EqualValues(t, types.NextGame{
		IdValue:  4,
		FifoHead: "1",
		FifoTail: "3",
	}, nextGame)
	game1, found1 := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "1")
	require.True(t, found1)
	require.EqualValues(t, types.StoredGame{
		Creator:   alice,
		Index:     "1",
		Game:      "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       bob,
		Black:     carol,
		MoveCount: uint64(0),
		BeforeId:  "-1",
		AfterId:   "3",
	}, game1)
	game3, found3 := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "3")
	require.True(t, found3)
	require.EqualValues(t, types.StoredGame{
		Creator:   carol,
		Index:     "3",
		Game:      "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*",
		Turn:      "b",
		Red:       alice,
		Black:     bob,
		MoveCount: uint64(0),
		BeforeId:  "1",
		AfterId:   "-1",
	}, game3)
}
