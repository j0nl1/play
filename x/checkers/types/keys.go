package types

import "time"

const (
	// ModuleName defines the module name
	ModuleName = "checkers"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_checkers"

	// StoredGame keys
	StoredGameEventKey     = "NewGameCreated"
	StoredGameEventCreator = "Creator"
	StoredGameEventIndex   = "Index"
	StoredGameEventRed     = "Red"
	StoredGameEventBlack   = "Black"

	// PlayMove keys
	PlayMoveEventKey       = "MovePlayed"
	PlayMoveEventCreator   = "Creator"
	PlayMoveEventIdValue   = "IdValue"
	PlayMoveEventCapturedX = "CapturedX"
	PlayMoveEventCapturedY = "CapturedY"
	PlayMoveEventWinner    = "Winner"

	// RejectGame keys
	RejectGameEventKey     = "GameRejected"
	RejectGameEventCreator = "Creator"
	RejectGameEventIdValue = "IdValue"

	// Deadline keys
	MaxTurnDuration = time.Duration(1 * 3_600 * 1000_000_000) // 1 hour
	DeadlineLayout  = "2006-01-02 15:04:05.999999999 +0000 UTC"

	// Fifo
	NoFifoIdKey = "-1"

	// Auto-expiring
	ForfeitGameEventKey     = "GameForfeited"
	ForfeitGameEventIdValue = "IdValue"
	ForfeitGameEventWinner  = "Winner"

	// Wager
	StoredGameEventWager = "Wager"
	StoredGameEventToken = "Token"

	// Gas
	CreateGameGas = 10
	PlayMoveGas   = 10
	RejectGameGas = 0
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

const (
	NextGameKey = "NextGame-value-"
)
