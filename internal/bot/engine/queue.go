package engine

import "github.com/ivan1993spb/snake-bot/internal/types"

type Position struct {
	dot      types.Dot
	score    int
	distance int
}

type Queue struct {
	positions []*Position
	enqueued  map[types.Dot]struct{}
}

const queueLength = 1024

func NewQueue() *Queue {
	return &Queue{
		positions: make([]*Position, 0, queueLength),
		enqueued:  make(map[types.Dot]struct{}),
	}
}

func (q *Queue) Len() int {
	return len(q.positions)
}

func (q *Queue) Less(i, j int) bool {
	if q.positions[i].score > q.positions[j].score {
		return true
	}
	if q.positions[i].score < q.positions[j].score {
		return false
	}
	return q.positions[i].distance < q.positions[j].distance

}

func (q *Queue) Swap(i, j int) {
	q.positions[i], q.positions[j] = q.positions[j], q.positions[i]
}

func (q *Queue) Push(x interface{}) {
	item := x.(*Position)
	q.positions = append(q.positions, item)
	q.enqueued[item.dot] = struct{}{}
}

func (q *Queue) Pop() interface{} {
	n := len(q.positions)
	item := q.positions[n-1]
	q.positions[n-1] = nil
	q.positions = q.positions[0 : n-1]
	delete(q.enqueued, item.dot)
	return item
}

func (q *Queue) Enqueued(dot types.Dot) bool {
	_, ok := q.enqueued[dot]
	return ok

}
