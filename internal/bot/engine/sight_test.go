package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ivan1993spb/snake-bot/internal/types"
)

func Test_NewSight(t *testing.T) {
	const format = "failed at test number %d"

	tests := []struct {
		area     Area
		pos      types.Dot
		distance uint8
		expected Sight
	}{
		{
			area: Area{
				Width:  100,
				Height: 50,
			},
			pos: types.Dot{
				X: 10,
				Y: 13,
			},
			distance: 3,
			expected: Sight{
				area: Area{
					Width:  100,
					Height: 50,
				},
				topLeft: types.Dot{
					X: 7,
					Y: 10,
				},
				zeroedBottomRight: types.Dot{
					X: 6,
					Y: 6,
				},
			},
		},
		{
			area: Area{
				Width:  3,
				Height: 3,
			},
			pos: types.Dot{
				X: 1,
				Y: 1,
			},
			distance: 1,
			expected: Sight{
				area: Area{
					Width:  3,
					Height: 3,
				},
				topLeft: types.Dot{
					X: 0,
					Y: 0,
				},
				zeroedBottomRight: types.Dot{
					X: 2,
					Y: 2,
				},
			},
		},
		{
			area: Area{
				Width:  4,
				Height: 4,
			},
			pos: types.Dot{
				X: 1,
				Y: 1,
			},
			distance: 2,
			expected: Sight{
				area: Area{
					Width:  4,
					Height: 4,
				},
				topLeft: types.Dot{
					X: 0,
					Y: 0,
				},
				zeroedBottomRight: types.Dot{
					X: 2,
					Y: 2,
				},
			},
		},
		{
			area: Area{
				Width:  10,
				Height: 4,
			},
			pos: types.Dot{
				X: 1,
				Y: 1,
			},
			distance: 2,
			expected: Sight{
				area: Area{
					Width:  10,
					Height: 4,
				},
				topLeft: types.Dot{
					X: 0,
					Y: 0,
				},
				zeroedBottomRight: types.Dot{
					X: 2,
					Y: 2,
				},
			},
		},
		{
			area: Area{
				Width:  13,
				Height: 32,
			},
			pos: types.Dot{
				X: 1,
				Y: 2,
			},
			distance: 5,
			expected: Sight{
				area: Area{
					Width:  13,
					Height: 32,
				},
				topLeft: types.Dot{
					X: 9,
					Y: 29,
				},
				zeroedBottomRight: types.Dot{
					X: 10,
					Y: 10,
				},
			},
		},
	}

	for i, test := range tests {
		sight := NewSight(test.area, test.pos, test.distance)
		assert.Equal(t, test.expected, sight, format, i+1)
	}
}

func Test_Sight_Absolute(t *testing.T) {
	const msgFormat = "%d: x%d y%d => %s instead of %s"

	tests := []struct {
		relX, relY uint8
		sight      Sight
		expect     types.Dot
	}{
		{
			relX: 0,
			relY: 0,
			sight: Sight{
				area: Area{
					Width:  13,
					Height: 32,
				},
				topLeft: types.Dot{
					X: 1,
					Y: 3,
				},
				zeroedBottomRight: types.Dot{
					X: 10,
					Y: 10,
				},
			},
			expect: types.Dot{
				X: 1,
				Y: 3,
			},
		},
		{
			relX: 4,
			relY: 5,
			sight: Sight{
				area: Area{
					Width:  20,
					Height: 20,
				},
				topLeft: types.Dot{
					X: 10,
					Y: 3,
				},
				zeroedBottomRight: types.Dot{
					X: 9,
					Y: 9,
				},
			},
			expect: types.Dot{
				X: 14,
				Y: 8,
			},
		},
		{
			relX: 7,
			relY: 1,
			sight: Sight{
				area: Area{
					Width:  20,
					Height: 20,
				},
				topLeft: types.Dot{
					X: 15,
					Y: 5,
				},
				zeroedBottomRight: types.Dot{
					X: 9,
					Y: 9,
				},
			},
			expect: types.Dot{
				X: 2,
				Y: 6,
			},
		},
		{
			relX: 68,
			relY: 89,
			sight: Sight{
				area: Area{
					Width:  255,
					Height: 255,
				},
				topLeft: types.Dot{
					X: 205,
					Y: 178,
				},
				zeroedBottomRight: types.Dot{
					X: 100,
					Y: 100,
				},
			},
			expect: types.Dot{
				X: 18,
				Y: 12,
			},
		},
	}

	for i, test := range tests {
		dot := test.sight.Absolute(test.relX, test.relY)
		assert.Equalf(t, test.expect, dot, msgFormat,
			i+1, test.relX, test.relY, dot, test.expect)
	}
}

func Test_Sight_Relative(t *testing.T) {
	// TODO: test this.
	t.Skip()
}

func Test_Sight_Relative_To_Absolute(t *testing.T) {
	s := Sight{
		area: Area{
			Width:  255,
			Height: 255,
		},
		topLeft: types.Dot{
			X: 250,
			Y: 249,
		},
		zeroedBottomRight: types.Dot{
			X: 100,
			Y: 100,
		},
	}

	for _, dot := range s.Dots() {
		x, y := s.Relative(dot)
		res := s.Absolute(x, y)
		assert.Equal(t, dot, res)
	}
}

func Test_Sight_Seen(t *testing.T) {
	const format = "failed at test number %d"

	tests := []struct {
		sight    Sight
		dot      types.Dot
		expected bool
	}{
		{
			sight: Sight{
				area: Area{
					Width:  13,
					Height: 32,
				},
				topLeft: types.Dot{
					X: 9,
					Y: 29,
				},
				zeroedBottomRight: types.Dot{
					X: 10,
					Y: 10,
				},
			},
			dot: types.Dot{
				X: 0,
				Y: 0,
			},
			expected: true,
		},
		{
			sight: Sight{
				area: Area{
					Width:  100,
					Height: 50,
				},
				topLeft: types.Dot{
					X: 7,
					Y: 10,
				},
				zeroedBottomRight: types.Dot{
					X: 6,
					Y: 6,
				},
			},
			dot: types.Dot{
				X: 9,
				Y: 10,
			},
			expected: true,
		},
		{
			sight: Sight{
				area: Area{
					Width:  13,
					Height: 32,
				},
				topLeft: types.Dot{
					X: 9,
					Y: 29,
				},
				zeroedBottomRight: types.Dot{
					X: 10,
					Y: 10,
				},
			},
			dot: types.Dot{
				X: 12,
				Y: 30,
			},
			expected: true,
		},
		{
			sight: Sight{
				area: Area{
					Width:  13,
					Height: 32,
				},
				topLeft: types.Dot{
					X: 9,
					Y: 29,
				},
				zeroedBottomRight: types.Dot{
					X: 10,
					Y: 10,
				},
			},
			dot: types.Dot{
				X: 21,
				Y: 2,
			},
			expected: false,
		},
	}

	for i, test := range tests {
		assert.Equal(t, test.expected, test.sight.Seen(test.dot), format, i+1)
	}
}

func Test_Sight_DotsDoesntRepeat(t *testing.T) {
	const msgFormat = "index %d; dot %s"

	const (
		width  = 255
		height = 255
	)

	check := map[types.Dot]struct{}{}

	a := NewArea(width, height)
	d := uint8(100)
	pos := types.Dot{
		X: 0,
		Y: 0,
	}
	s := NewSight(a, pos, d)

	for i, dot := range s.Dots() {
		if _, ok := check[dot]; ok {
			assert.Falsef(t, ok, msgFormat, i, dot)
		}
		check[dot] = struct{}{}
	}
}
