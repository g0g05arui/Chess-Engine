package computed

import (
	"sync"

	engine "github.com/g0g05arui/chess-engine/game_state"
)

type CacheKey = struct {
	Fen       string
	WhiteTurn bool
}

type CacheValue = struct {
	BestMove engine.Move
	Depth    int
}

type DeepCacheKey struct {
	Fen       string
	WhiteTurn bool
	Depth     int
}

type DeepCacheValue struct {
	BestMove engine.Move
	Depth    int
}

var DeepCache = make(map[DeepCacheKey]DeepCacheValue)
var Cache = make(map[CacheKey]engine.Move)
var InProgress sync.Map // key = DeepCacheKey, value = struct{}{}
