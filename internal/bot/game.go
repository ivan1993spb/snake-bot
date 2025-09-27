package bot

import (
	"sync"

	"github.com/ivan1993spb/snake-bot/internal/bot/engine"
	"github.com/ivan1993spb/snake-bot/internal/types"
)

type gameState uint8

const (
	gameStateWait gameState = iota
	gameStateReady
)

type Game struct {
	state   gameState
	mux     *sync.RWMutex
	objects map[uint32]*types.Object
	area    engine.Area
	_map    *engine.Map
}

func NewGame() *Game {
	return &Game{
		state:   gameStateWait,
		mux:     &sync.RWMutex{},
		objects: make(map[uint32]*types.Object),
	}
}

func (g *Game) Create(object *types.Object) {
	g.mux.Lock()
	defer g.mux.Unlock()
	g.objects[object.Id] = object
	if g.state == gameStateReady {
		g._map.SaveObject(object)
	}
}

func (g *Game) Update(object *types.Object) {
	g.mux.Lock()
	defer g.mux.Unlock()
	old := g.objects[object.Id]
	g.objects[object.Id] = object
	if g.state == gameStateReady {
		if old != nil {
			g._map.Clear(old.GetDots())
		}
		g._map.SaveObject(object)
	}
}

func (g *Game) UpdateV2(update *types.UpdateV2) {
	g.mux.Lock()
	defer g.mux.Unlock()
	object, ok := g.objects[update.Id]
	if !ok {
		return
	}
	// Only for snakes because only snakes grow.
	if object.Type == types.ObjectTypeSnake && update.Add != nil {
		// Add new head to the snake's body
		object.Dots = append([]types.Dot{*update.Add}, object.Dots...)
		if g.state == gameStateReady {
			g._map.SaveObjectDot(object, *update.Add)
		}
	} else if object.Type == types.ObjectTypeMouse && update.Dot != nil {
		// Only for mice because only mice changes its coordinates.
		if g.state == gameStateReady {
			g._map.ClearDot(object.Dot)
			g._map.SaveObjectDot(object, *update.Dot)
		}
		object.Dot = *update.Dot
	}
	if update.Del != nil {
		// Most of updates are coming for snakes. So we search from the end
		// because we are looking to delete the tail. The same logic
		// applies for other food objects.
		for i := len(object.Dots) - 1; i >= 0; i-- {
			if object.Dots[i] == *update.Del {
				// Now tails is found.
				object.Dots = append(object.Dots[:i], object.Dots[i+1:]...)
				if g.state == gameStateReady {
					g._map.ClearDot(*update.Del)
				}
				break
			}
		}
	}
}

func (g *Game) Delete(object *types.Object) {
	g.mux.Lock()
	defer g.mux.Unlock()
	delete(g.objects, object.Id)
	if g.state == gameStateReady {
		g._map.Clear(object.GetDots())
	}
}

func (g *Game) GetObject(id uint32) (*types.Object, bool) {
	g.mux.RLock()
	defer g.mux.RUnlock()
	if object, ok := g.objects[id]; ok {
		return object, true
	}
	return nil, false
}

func (g *Game) LookAround(sight engine.Sight) *engine.HashmapSight {
	g.mux.RLock()
	defer g.mux.RUnlock()

	if g.state == gameStateReady {
		return g._map.LookAround(sight)
	}

	return nil
}

func (g *Game) Size(width, height uint8) {
	g.mux.Lock()
	defer g.mux.Unlock()
	g.area = engine.Area{
		Width:  width,
		Height: height,
	}
	g._map = engine.NewMap(g.area)
	g.state = gameStateReady
}

func (g *Game) GetArea() engine.Area {
	g.mux.RLock()
	defer g.mux.RUnlock()
	return g.area
}
