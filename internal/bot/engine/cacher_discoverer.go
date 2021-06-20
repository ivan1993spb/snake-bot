package engine

import (
	"github.com/ivan1993spb/snake-bot/internal/types"
)

// TODO: Consider using cache size.

var _ Discoverer = (*CacherDiscoverer)(nil)

type CacherDiscoverer struct {
	discoverer Discoverer

	path   []types.Dot
	scores []int

	expected int
}

func NewCacherDiscoverer(discoverer Discoverer) Discoverer {
	return &CacherDiscoverer{
		discoverer: discoverer,
	}
}

func (d *CacherDiscoverer) Discover(head types.Dot, area Area,
	sight Sight, scores *HashmapSight) []types.Dot {

	d.update(head)

	if d.expired(scores) {
		d.path = d.discoverer.Discover(head, area, sight, scores)
		d.score(scores)
	}

	return d.path
}

func (d *CacherDiscoverer) update(head types.Dot) {
	if len(d.path) == 0 {
		return
	}
	if i := d.index(head); i > -1 {
		d.path = d.path[i+1:]
		d.scores = d.scores[i+1:]
	}
}

func (d *CacherDiscoverer) index(head types.Dot) int {
	for i, dot := range d.path {
		if dot == head {
			return i
		}
	}
	return -1
}

const visibilityThreshold = 10

func (d *CacherDiscoverer) expired(scores *HashmapSight) bool {
	if len(d.path) < visibilityThreshold || d.expected <= 0 {
		return true
	}
	for i, dot := range d.path {
		score, _ := scores.AccessDefault(dot, 0).(int)
		if d.scores[i] != score {
			return true
		}
	}
	return false
}

func (d *CacherDiscoverer) score(scores *HashmapSight) {
	d.expected = 0
	d.scores = make([]int, len(d.path))

	for i, dot := range d.path {
		score, _ := scores.AccessDefault(dot, 0).(int)
		d.scores[i] = score
		d.expected += d.scores[i]
	}
}
