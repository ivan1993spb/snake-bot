package engine

import "github.com/ivan1993spb/snake-bot/internal/types"

type Discoverer interface {
	Discover(head types.Dot, area Area, sight Sight,
		scores *HashmapSight) []types.Dot
}
