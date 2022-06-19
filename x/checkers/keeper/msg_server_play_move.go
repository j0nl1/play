package keeper

import (
	"context"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/j0nl1/play/x/checkers/rules"
	"github.com/j0nl1/play/x/checkers/types"
)

func (k msgServer) PlayMove(goCtx context.Context, msg *types.MsgPlayMove) (*types.MsgPlayMoveResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	storedGame, foundGame := k.Keeper.GetStoredGame(ctx, msg.IdValue)

	if !foundGame {
		return nil, types.ErrGameNotFound
	}

	if storedGame.Winner != rules.PieceStrings[rules.NO_PLAYER] {
		return nil, types.ErrGameFinished
	}

	if storedGame.MoveCount <= 1 {
		err := k.Keeper.CollectWager(ctx, &storedGame)
		if err != nil {
			return nil, err
		}
	}

	isRed := strings.Compare(storedGame.Red, msg.Creator) == 0
	isBlack := strings.Compare(storedGame.Black, msg.Creator) == 0

	var player rules.Player

	if !isRed && !isBlack {
		return nil, types.ErrCreatorNotPlayer
	} else if isRed && isBlack {
		player = rules.StringPieces[storedGame.Turn].Player
	} else if isRed {
		player = rules.RED_PLAYER
	} else {
		player = rules.BLACK_PLAYER
	}

	game, parsedError := storedGame.ParseGame()

	if parsedError != nil {
		panic(parsedError.Error())
	}

	if !game.TurnIs(player) {
		return nil, types.ErrNotPlayerTurn
	}

	captured, moveErr := game.Move(
		rules.Pos{
			X: int(msg.FromX),
			Y: int(msg.FromY),
		},
		rules.Pos{
			X: int(msg.ToX),
			Y: int(msg.ToY),
		},
	)
	if moveErr != nil {
		return nil, sdkerrors.Wrapf(types.ErrWrongMove, moveErr.Error())
	}

	storedGame.MoveCount++
	storedGame.Deadline = types.FormatDeadline(types.GetNextDeadline(ctx))
	storedGame.Winner = rules.PieceStrings[game.Winner()]

	nextGame, found := k.Keeper.GetNextGame(ctx)

	if !found {
		panic("NextGame not found")
	}
	if storedGame.Winner == rules.PieceStrings[rules.NO_PLAYER] {
		k.Keeper.SendToFifoTail(ctx, &storedGame, &nextGame)
	} else {
		k.Keeper.RemoveFromFifo(ctx, &storedGame, &nextGame)
		k.Keeper.MustPayWinnings(ctx, &storedGame)
	}

	storedGame.Game = game.String()
	storedGame.Turn = rules.PieceStrings[game.Turn]

	k.Keeper.SetStoredGame(ctx, storedGame)
	k.Keeper.SetNextGame(ctx, nextGame)

	ctx.GasMeter().ConsumeGas(types.PlayMoveGas, "Play a move")

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "checkers"),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.PlayMoveEventKey),
			sdk.NewAttribute(types.PlayMoveEventCreator, msg.Creator),
			sdk.NewAttribute(types.PlayMoveEventIdValue, msg.IdValue),
			sdk.NewAttribute(types.PlayMoveEventCapturedX, strconv.FormatInt(int64(captured.X), 10)),
			sdk.NewAttribute(types.PlayMoveEventCapturedY, strconv.FormatInt(int64(captured.Y), 10)),
			sdk.NewAttribute(types.PlayMoveEventWinner, game.Winner().Color),
		),
	)

	return &types.MsgPlayMoveResponse{
		IdValue:   msg.IdValue,
		CapturedX: uint64(captured.X),
		CapturedY: uint64(captured.Y),
		Winner:    game.Winner().Color,
	}, nil
}
