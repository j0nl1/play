package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/j0nl1/play/x/checkers/rules"
	"github.com/j0nl1/play/x/checkers/types"
)

func (k msgServer) RejectGame(goCtx context.Context, msg *types.MsgRejectGame) (*types.MsgRejectGameResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	storedGame, found := k.Keeper.GetStoredGame(ctx, msg.IdValue)
	if !found {
		return nil, sdkerrors.Wrapf(types.ErrGameNotFound, "game not found %s", msg.IdValue)
	}

	if storedGame.Winner != rules.PieceStrings[rules.NO_PLAYER] {
		return nil, types.ErrGameFinished
	}

	if strings.Compare(storedGame.Red, msg.Creator) == 0 {
		if 1 < storedGame.MoveCount {
			return nil, types.ErrRedAlreadyPlayed
		}
	} else if strings.Compare(storedGame.Black, msg.Creator) == 0 {
		if 0 < storedGame.MoveCount {
			return nil, types.ErrBlackAlreadyPlayed
		}
	} else {
		return nil, types.ErrCreatorNotPlayer
	}

	nextGame, found := k.Keeper.GetNextGame(ctx)
	if !found {
		panic("NextGame not found")
	}
	k.Keeper.RemoveFromFifo(ctx, &storedGame, &nextGame)
	k.Keeper.RemoveStoredGame(ctx, msg.IdValue)
	k.Keeper.MustRefundWager(ctx, &storedGame)

	k.Keeper.SetNextGame(ctx, nextGame)

	ctx.GasMeter().ConsumeGas(types.RejectGameGas, "Reject game")

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, "checkers"),
			sdk.NewAttribute(sdk.AttributeKeyAction, types.RejectGameEventKey),
			sdk.NewAttribute(types.RejectGameEventCreator, msg.Creator),
			sdk.NewAttribute(types.RejectGameEventIdValue, msg.IdValue),
		),
	)

	return &types.MsgRejectGameResponse{}, nil
}
