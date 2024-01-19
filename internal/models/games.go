package models

import "sort"

type Games struct {
	Games []*Game `json:"games" yaml:"games"`
}

func NewGames(state map[int]int) *Games {
	g := &Games{}
	g.Games = make([]*Game, 0, len(state))
	for game, bots := range state {
		if bots > 0 {
			g.Games = append(g.Games, &Game{
				Game: game,
				Bots: bots,
			})
		}
	}

	sort.Slice(g.Games, func(i, j int) bool {
		return g.Games[i].Game < g.Games[j].Game
	})

	return g
}

func (g *Games) ToMapState() map[int]int {
	state := make(map[int]int, len(g.Games))
	for _, game := range g.Games {
		state[game.Game] = game.Bots
	}
	return state
}
