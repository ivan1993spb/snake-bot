package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ivan1993spb/snake-bot/internal/types"
)

func Test_NewArea(t *testing.T) {
	const (
		w = 132
		h = 43
	)
	a := NewArea(w, h)
	assert.Equal(t, uint8(w), a.Width)
	assert.Equal(t, uint8(h), a.Height)
}

func Test_Area_Navigate(t *testing.T) {
	const msgFormat = "%d: %s.Navigate(%s) = %s"

	tests := []struct {
		area   Area
		input  types.Dot
		output []types.Dot
	}{
		{
			area: Area{
				Width:  100,
				Height: 101,
			},
			input: types.Dot{
				X: 0,
				Y: 0,
			},
			output: []types.Dot{
				{
					X: 0,
					Y: 100,
				},
				{
					X: 0,
					Y: 1,
				},
				{
					X: 1,
					Y: 0,
				},
				{
					X: 99,
					Y: 0,
				},
			},
		},
		{
			area: Area{
				Width:  100,
				Height: 101,
			},
			input: types.Dot{
				X: 99,
				Y: 0,
			},
			output: []types.Dot{
				{
					X: 99,
					Y: 100,
				},
				{
					X: 99,
					Y: 1,
				},
				{
					X: 0,
					Y: 0,
				},
				{
					X: 98,
					Y: 0,
				},
			},
		},
		{
			area: Area{
				Width:  100,
				Height: 101,
			},
			input: types.Dot{
				X: 99,
				Y: 100,
			},
			output: []types.Dot{
				{
					X: 99,
					Y: 99,
				},
				{
					X: 99,
					Y: 0,
				},
				{
					X: 0,
					Y: 100,
				},
				{
					X: 98,
					Y: 100,
				},
			},
		},
		{
			area: Area{
				Width:  55,
				Height: 24,
			},
			input: types.Dot{
				X: 0,
				Y: 23,
			},
			output: []types.Dot{
				{
					X: 0,
					Y: 22,
				},
				{
					X: 0,
					Y: 0,
				},
				{
					X: 1,
					Y: 23,
				},
				{
					X: 54,
					Y: 23,
				},
			},
		},
		{
			area: Area{
				Width:  55,
				Height: 24,
			},
			input: types.Dot{
				X: 12,
				Y: 4,
			},
			output: []types.Dot{
				{
					X: 12,
					Y: 3,
				},
				{
					X: 12,
					Y: 5,
				},
				{
					X: 13,
					Y: 4,
				},
				{
					X: 11,
					Y: 4,
				},
			},
		},
	}

	for i, test := range tests {
		dots := test.area.Navigate(test.input)
		assert.Equalf(t, test.output, dots, msgFormat,
			i+1, test.area, test.input, dots)
	}
}

func Test_Area_FindDirection(t *testing.T) {
	const msgFormat = "%d: %s.FindDirection(%s, %s) = %s"

	tests := []struct {
		area   Area
		from   types.Dot
		to     types.Dot
		expect types.Direction
	}{
		{
			area: Area{
				Width:  10,
				Height: 11,
			},
			from:   types.Dot{},
			to:     types.Dot{},
			expect: types.DirectionNorth,
		},
		{
			area: Area{
				Width:  10,
				Height: 11,
			},
			from: types.Dot{},
			to: types.Dot{
				X: 12,
				Y: 23,
			},
			expect: types.DirectionNorth,
		},
		{
			area: Area{
				Width:  10,
				Height: 11,
			},
			from: types.Dot{
				X: 12,
				Y: 23,
			},
			to:     types.Dot{},
			expect: types.DirectionNorth,
		},

		// TODO: Add cases.
	}

	for i, test := range tests {
		dir := test.area.FindDirection(test.from, test.to)
		assert.Equalf(t, test.expect, dir, msgFormat,
			i+1, test.area, test.from, test.to, dir)
	}
}
