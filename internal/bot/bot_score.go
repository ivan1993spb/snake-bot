package bot

import (
	"github.com/ivan1993spb/snake-bot/internal/bot/engine"
	"github.com/ivan1993spb/snake-bot/internal/types"
)

type scoreType uint8

const (
	scoreTypeCollapse scoreType = iota
	scoreTypeFoodApple
	scoreTypeFoodCorpse
	scoreTypeFoodWatermelon
	scoreTypeFoodMouse
	scoreTypeHunt
)

var defaultBehavior = map[scoreType]int{
	scoreTypeCollapse:       -1000,
	scoreTypeFoodApple:      1,
	scoreTypeFoodCorpse:     2,
	scoreTypeFoodWatermelon: 5,
	scoreTypeFoodMouse:      15,
	scoreTypeHunt:           30,
}

func (b *Bot) score(objects *engine.HashmapSight) *engine.HashmapSight {
	// TODO: Estimate the real prey length.
	const preyLen = 3

	scores := objects.Reflect()

	objects.ForEach(func(dot types.Dot, v interface{}) {
		object, ok := v.(*types.Object)
		if !ok {
			return
		}
		if object.Id == b.myId {
			scores.Assign(dot, defaultBehavior[scoreTypeCollapse])
		} else if object.Type == types.ObjectTypeSnake {
			if len(object.Dots) == preyLen {
				scores.Assign(dot, defaultBehavior[scoreTypeHunt])
			} else {
				scores.Assign(dot, defaultBehavior[scoreTypeCollapse])
			}
		} else if object.Type == types.ObjectTypeWall {
			scores.Assign(dot, defaultBehavior[scoreTypeCollapse])
		} else {
			switch object.Type {
			case types.ObjectTypeApple:
				scores.Assign(dot, defaultBehavior[scoreTypeFoodApple])
			case types.ObjectTypeCorpse:
				scores.Assign(dot, defaultBehavior[scoreTypeFoodCorpse])
			case types.ObjectTypeWatermelon:
				scores.Assign(dot, defaultBehavior[scoreTypeFoodWatermelon])
			case types.ObjectTypeMouse:
				scores.Assign(dot, defaultBehavior[scoreTypeFoodMouse])
			default:
				// Avoid unknown objects.
				scores.Assign(dot, defaultBehavior[scoreTypeCollapse])
			}
		}
	})

	return scores
}
