package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/j0nl1/play/x/checkers/types"
)

func (k Keeper) RemoveFromFifo(ctx sdk.Context, game *types.StoredGame, info *types.NextGame) {
	if game.BeforeId != types.NoFifoIdKey {
		beforeElement, found := k.GetStoredGame(ctx, game.BeforeId)
		if !found {
			panic("Element before in Fifo was not found")
		}
		beforeElement.AfterId = game.AfterId
		k.SetStoredGame(ctx, beforeElement)
		if game.AfterId == types.NoFifoIdKey {
			info.FifoTail = beforeElement.Index
		}

	} else if info.FifoHead == game.Index {
		info.FifoHead = game.AfterId
	}

	if game.AfterId != types.NoFifoIdKey {
		afterElement, found := k.GetStoredGame(ctx, game.AfterId)
		if !found {
			panic("Element after in Fifo was not found")
		}
		afterElement.BeforeId = game.BeforeId
		k.SetStoredGame(ctx, afterElement)
		if game.BeforeId == types.NoFifoIdKey {
			info.FifoHead = afterElement.Index
		}

	} else if info.FifoTail == game.Index {
		info.FifoTail = game.BeforeId
	}
	game.BeforeId = types.NoFifoIdKey
	game.AfterId = types.NoFifoIdKey
}

func (k Keeper) SendToFifoTail(ctx sdk.Context, game *types.StoredGame, info *types.NextGame) {
	if info.FifoHead == types.NoFifoIdKey && info.FifoTail == types.NoFifoIdKey {
		game.BeforeId = types.NoFifoIdKey
		game.AfterId = types.NoFifoIdKey
		info.FifoHead = game.Index
		info.FifoTail = game.Index
	} else if info.FifoHead == types.NoFifoIdKey || info.FifoTail == types.NoFifoIdKey {
		panic("Fifo should have both head and tail or none")
	} else if info.FifoTail == game.Index {
	} else {

		k.RemoveFromFifo(ctx, game, info)

		currentTail, found := k.GetStoredGame(ctx, info.FifoTail)
		if !found {
			panic("Current Fifo tail was not found")
		}
		currentTail.AfterId = game.Index
		k.SetStoredGame(ctx, currentTail)

		game.BeforeId = currentTail.Index
		info.FifoTail = game.Index
	}
}
