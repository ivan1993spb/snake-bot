package engine

import "github.com/ivan1993spb/snake-bot/internal/types"

func backtrack(from, to types.Dot, distance int,
	prev map[types.Dot]types.Dot) []types.Dot {
	path := make([]types.Dot, distance)
	current := to
	for from != current {
		distance--
		path[distance] = current
		current = prev[current]
	}
	return path
}
