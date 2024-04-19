package bot

import (
	"context"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ivan1993spb/snake-bot/internal/bot/engine"
	"github.com/ivan1993spb/snake-bot/internal/types"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

const botTickTime = time.Millisecond * 200

const lookupDistance uint8 = 50

const directionExpireTime = time.Millisecond * 10

const (
	stateWait = iota
	stateExplore
)

type World interface {
	LookAround(sight engine.Sight) *engine.HashmapSight
	GetArea() engine.Area
	GetObject(id uint32) (*types.Object, bool)
}

type Bot struct {
	state uint32

	myId  uint32
	world World

	discoverer engine.Discoverer

	lastPosition  types.Dot
	lastDirection types.Direction
}

func NewBot(world World, discoverer engine.Discoverer) *Bot {
	return &Bot{
		world:      world,
		discoverer: discoverer,
	}
}

func NewDijkstrasBot(world World) *Bot {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	discoverer := engine.NewDijkstrasDiscoverer(r)
	cacher := engine.NewCacherDiscoverer(discoverer)
	return NewBot(world, cacher)
}

func (b *Bot) Run(ctx context.Context) <-chan types.Direction {
	chout := make(chan types.Direction)

	go func() {
		defer close(chout)
		ticker := time.NewTicker(botTickTime)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				direction, ok := b.operate(ctx)
				if ok {
					ctx, cancel := context.WithTimeout(ctx, directionExpireTime)
					select {
					case chout <- direction:
					case <-ctx.Done():
					}
					cancel()
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return chout
}

func (b *Bot) getState() uint32 {
	return atomic.LoadUint32(&b.state)
}

func (b *Bot) getMe() uint32 {
	return atomic.LoadUint32(&b.myId)
}

func (b *Bot) operate(ctx context.Context) (types.Direction, bool) {
	if b.getState() != stateExplore {
		return types.DirectionZero, false
	}

	me, ok := b.world.GetObject(b.getMe())
	if !ok {
		return types.DirectionZero, false
	}

	area := b.world.GetArea()
	objectDots := me.GetDots()
	if len(objectDots) == 0 {
		return types.DirectionZero, false
	}
	head := objectDots[0]
	sight := engine.NewSight(area, head, lookupDistance)

	objects := b.world.LookAround(sight)
	scores := b.score(objects)

	path := b.discoverer.Discover(head, area, sight, scores)
	if len(path) == 0 {
		return types.DirectionZero, false
	}
	direction := area.FindDirection(head, path[0])

	if area.FindDirection(objectDots[1], objectDots[0]) == direction {
		// the same direction
		return types.DirectionZero, false
	}
	if b.lastDirection != direction || b.lastPosition != head {
		utils.GetLogger(ctx).WithFields(logrus.Fields{
			"head":      head,
			"direction": direction,
		}).Debugln("change direction")
		b.lastDirection = direction
		b.lastPosition = head

		return direction, true
	}
	return types.DirectionZero, false
}

func (b *Bot) Countdown(sec int) {
	// TODO: Shut down the snake if the stateWait
	atomic.StoreUint32(&b.state, stateWait)
}

func (b *Bot) Me(id uint32) {
	atomic.StoreUint32(&b.myId, id)
	atomic.StoreUint32(&b.state, stateExplore)
}
