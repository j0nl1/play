package keeper_test

import (
	"context"
	"math/rand"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/j0nl1/play/testutil/keeper"
	"github.com/j0nl1/play/x/checkers"
	"github.com/j0nl1/play/x/checkers/keeper"
	"github.com/j0nl1/play/x/checkers/rules"
	"github.com/j0nl1/play/x/checkers/types"
	"github.com/stretchr/testify/require"
)

const (
	alice = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d3"
	bob   = "cosmos1xyxs3skf3f4jfqeuv89yyaqvjc6lffavxqhc8g"
	carol = "cosmos1e0w5t53nrq7p66fye6c8p0ynyhf6y24l4yuxd7"
)

func setupMsgServerCreateGame(t testing.TB) (types.MsgServer, keeper.Keeper, context.Context) {
	k, ctx := keepertest.CheckersKeeper(t)
	checkers.InitGenesis(ctx, *k, *types.DefaultGenesis())
	return keeper.NewMsgServerImpl(*k), *k, sdk.WrapSDKContext(ctx)
}

func createMoreThanOneGame(context context.Context, msgServer types.MsgServer, numberOfGames uint8) {
	for i := uint8(0); i < numberOfGames; i++ {
		msgServer.CreateGame(context, &types.MsgCreateGame{
			Creator: alice,
			Red:     bob,
			Black:   carol,
		})
	}
}

func TestCreateGame(t *testing.T) {
	msgServer, _, context := setupMsgServerCreateGame(t)
	createResponse, err := msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Red:     bob,
		Black:   carol,
	})
	require.Nil(t, err)
	require.EqualValues(t, types.MsgCreateGameResponse{
		IdValue: "1",
	}, *createResponse)
}

func TestCreate1GameHasSaved(t *testing.T) {
	msgServer, keeper, context := setupMsgServerCreateGame(t)

	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Red:     bob,
		Black:   carol,
	})

	nextGame, nextGameFound := keeper.GetNextGame(sdk.UnwrapSDKContext(context))
	require.True(t, nextGameFound)
	require.EqualValues(t, types.NextGame{
		IdValue: 2,
	}, nextGame)

	storedGame, storedGameFound := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "1")
	require.True(t, storedGameFound)
	newGame := rules.New()
	require.EqualValues(t, types.StoredGame{
		Creator: alice,
		Index:   "1",
		Game:    newGame.String(),
		Turn:    rules.PieceStrings[newGame.Turn],
		Red:     bob,
		Black:   carol,
	}, storedGame)
}

func TestCreateXGames(t *testing.T) {
	msgServer, keeper, context := setupMsgServerCreateGame(t)

	numberOfGames := rand.Intn(30) + 1
	strNumberOfGames := strconv.FormatUint(uint64(numberOfGames), 10)

	createMoreThanOneGame(context, msgServer, uint8(numberOfGames))

	firstGame, firstGameFound := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), "1")
	// In case number of games is 1, lastGame will be equal to firstGame
	lastGame, lastGameFound := keeper.GetStoredGame(sdk.UnwrapSDKContext(context), strNumberOfGames)
	require.True(t, firstGameFound)
	require.True(t, lastGameFound)
	require.EqualValues(t, firstGame.Index, "1")
	require.EqualValues(t, lastGame.Index, strNumberOfGames)

}

func TestCreateXGamesGetAll(t *testing.T) {
	msgServer, keeper, context := setupMsgServerCreateGame(t)

	numberOfGames := rand.Intn(30) + 1

	createMoreThanOneGame(context, msgServer, uint8(numberOfGames))

	games := keeper.GetAllStoredGame(sdk.UnwrapSDKContext(context))

	require.EqualValues(t, numberOfGames, len(games))
}

func TestCreateGameFarFuture(t *testing.T) {
	msgSrvr, keeper, context := setupMsgServerCreateGame(t)
	keeper.SetNextGame(sdk.UnwrapSDKContext(context), types.NextGame{
		IdValue: 1024,
	})

	createResponse, err := msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Red:     bob,
		Black:   carol,
	})

	require.Nil(t, err)
	require.EqualValues(t, types.MsgCreateGameResponse{
		IdValue: "1024",
	}, *createResponse)
}

func TestCreate1GameEmitted(t *testing.T) {
	msgSrvr, _, context := setupMsgServerCreateGame(t)
	msgSrvr.CreateGame(context, &types.MsgCreateGame{
		Creator: alice,
		Red:     bob,
		Black:   carol,
	})
	ctx := sdk.UnwrapSDKContext(context)
	require.NotNil(t, ctx)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 1)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "message",
		Attributes: []sdk.Attribute{
			{Key: "module", Value: "checkers"},
			{Key: "action", Value: "NewGameCreated"},
			{Key: "Creator", Value: alice},
			{Key: "Index", Value: "1"},
			{Key: "Red", Value: bob},
			{Key: "Black", Value: carol},
		},
	}, event)
}
