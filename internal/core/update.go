package core

type update interface {
	diff(map[int]int) map[int]int
}

type updateOne struct {
	game int
	bots int
}

func (u *updateOne) diff(state map[int]int) map[int]int {
	return diffOne(u.game, u.bots, state)
}

var _ update = (*updateOne)(nil)

type updateBulk struct {
	state map[int]int
}

func (u *updateBulk) diff(state map[int]int) map[int]int {
	return diff(state, u.state)
}

var _ update = (*updateBulk)(nil)
