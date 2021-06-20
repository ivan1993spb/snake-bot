package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ivan1993spb/snake-bot/internal/types"
)

func TestCacherDiscoverer_update_ReturnsIfPathEmpty(t *testing.T) {
	d := &CacherDiscoverer{
		path:   []types.Dot{},
		scores: []int{},
	}

	head := types.Dot{1, 3}

	d.update(head)
}

func TestCacherDiscoverer_update_CutsPath(t *testing.T) {
	d := &CacherDiscoverer{
		path: []types.Dot{
			{1, 0},
			{1, 1},
			{1, 2},
			{1, 3},
			{1, 4},
			{1, 5},
		},
		scores: []int{
			1,
			0,
			3,
			1,
			0,
			0,
		},
	}

	head := types.Dot{1, 3}

	d.update(head)

	expectPath := []types.Dot{
		{1, 4},
		{1, 5},
	}
	expectScores := []int{
		0,
		0,
	}

	assert.Equal(t, expectPath, d.path)
	assert.Equal(t, expectScores, d.scores)
}

func TestCacherDiscoverer_update_RemovesCache(t *testing.T) {
	d := &CacherDiscoverer{
		path: []types.Dot{
			{1, 0},
			{1, 1},
			{1, 2},
			{1, 3},
			{1, 4},
			{1, 5},
		},
		scores: []int{
			1,
			0,
			3,
			1,
			0,
			0,
		},
	}

	head := types.Dot{1, 5}

	d.update(head)

	assert.Empty(t, d.path)
	assert.Empty(t, d.scores)
}
