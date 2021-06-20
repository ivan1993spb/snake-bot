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
