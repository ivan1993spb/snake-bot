package engine

import (
	"container/heap"
	"math/rand"

	"github.com/ivan1993spb/snake-bot/internal/types"
)

const maxPathLength = 50

var _ Discoverer = (*DijkstrasDiscoverer)(nil)

// DijkstrasDiscoverer discoves paths. It is based on
// Dijkstra's algorithm.
//
// Link: https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm
type DijkstrasDiscoverer struct {
	random *rand.Rand
}

func NewDijkstrasDiscoverer(r *rand.Rand) Discoverer {
	return &DijkstrasDiscoverer{
		random: r,
	}
}

func (d *DijkstrasDiscoverer) Discover(head types.Dot, area Area,
	sight Sight, scores *HashmapSight) []types.Dot {
	visited := make(map[types.Dot]struct{})
	prev := make(map[types.Dot]types.Dot)

	// Start with a negative score! In this case the algorithm
	// won't conclude that the current position is the best
	// even if all the observed dots are empty.
	best := &Position{
		dot:      head,
		score:    -1,
		distance: 0,
	}
	queue := NewQueue()
	queue.Push(best)
	heap.Init(queue)

	for queue.Len() > 0 {
		position := heap.Pop(queue).(*Position)
		dots := area.Navigate(position.dot)
		d.shuffle(dots)

		for _, dot := range dots {
			if !sight.Seen(dot) || queue.Enqueued(dot) {
				continue
			}
			if _, ok := visited[dot]; ok {
				continue
			}
			score, _ := scores.AccessDefault(dot, 0).(int)
			if score >= 0 {
				prev[dot] = position.dot
				distance := position.distance + 1

				if distance < maxPathLength {
					heap.Push(queue, &Position{
						dot:      dot,
						score:    position.score + score,
						distance: distance,
					})
				}
			}
		}

		visited[position.dot] = struct{}{}
		if best.score < position.score {
			best = position
		}
	}

	path := backtrack(head, best.dot, best.distance, prev)
	return path
}

func (d *DijkstrasDiscoverer) shuffle(dots []types.Dot) {
	d.random.Shuffle(len(dots), func(i, j int) {
		dots[i], dots[j] = dots[j], dots[i]
	})
}
