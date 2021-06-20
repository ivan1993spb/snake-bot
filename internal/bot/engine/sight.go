package engine

import (
	"github.com/ivan1993spb/snake-bot/internal/types"
)

type Sight struct {
	area    Area
	topLeft types.Dot

	zeroedBottomRight types.Dot

	width  uint8
	height uint8
}

// sightDivisor defines how many intervals of a given length
// we need to be able to fit within an area. In the 2D space
// the divisor has to be 2:
//                         ___
//                          ^
//                          |
//                          |
//                      distance
//                          |
//                          |
//                          v
//        |<---distance--->GAP<---distance--->|
//                          ^
//                          |
//                          |
//                      distance
//                          |
//                          |
//                          v
//                         ---
const sightDivisor = 2

const sightGap = 1

func NewSight(a Area, pos types.Dot, distance uint8) Sight {
	distance = a.FitDistance(distance, sightDivisor, sightGap)

	topLeft := types.Dot{
		X: pos.X - distance,
		Y: pos.Y - distance,
	}
	if pos.X < distance {
		topLeft.X += a.Width
	}
	if pos.Y < distance {
		topLeft.Y += a.Height
	}

	bottomRight := types.Dot{
		X: (pos.X + distance) % a.Width,
		Y: (pos.Y + distance) % a.Height,
	}

	zeroedBottomRight := types.Dot{
		X: bottomRight.X - topLeft.X,
		Y: bottomRight.Y - topLeft.Y,
	}
	if bottomRight.X < topLeft.X {
		zeroedBottomRight.X += a.Width
	}
	if bottomRight.Y < topLeft.Y {
		zeroedBottomRight.Y += a.Height
	}

	return Sight{
		area:    a,
		topLeft: topLeft,

		zeroedBottomRight: zeroedBottomRight,

		// TODO: Check on overflow.
		//width:  zeroedBottomRight.X + 1,
		//height: zeroedBottomRight.Y + 1,
	}
}

func (s Sight) Absolute(relX, relY uint8) types.Dot {
	x := uint16(s.topLeft.X) + uint16(relX)
	y := uint16(s.topLeft.Y) + uint16(relY)
	return types.Dot{
		X: uint8(x % uint16(s.area.Width)),
		Y: uint8(y % uint16(s.area.Height)),
	}
}

func (s Sight) Relative(dot types.Dot) (uint8, uint8) {
	x := dot.X - s.topLeft.X
	if s.topLeft.X > dot.X {
		x += s.area.Width
	}
	y := dot.Y - s.topLeft.Y
	if s.topLeft.Y > dot.Y {
		y += s.area.Height
	}
	return x, y
}

func (s Sight) Seen(dot types.Dot) bool {
	if zeroedX := dot.X - s.topLeft.X; dot.X >= s.topLeft.X {
		if zeroedX > s.zeroedBottomRight.X {
			return false
		}
	} else if zeroedX+s.area.Width > s.zeroedBottomRight.X {
		return false
	}
	if zeroedY := dot.Y - s.topLeft.Y; dot.Y >= s.topLeft.Y {
		if zeroedY > s.zeroedBottomRight.Y {
			return false
		}
	} else if zeroedY+s.area.Height > s.zeroedBottomRight.Y {
		return false
	}
	return true
}

func (s Sight) Dots() []types.Dot {
	w := (int(s.zeroedBottomRight.X) + 1)
	h := (int(s.zeroedBottomRight.Y) + 1)

	dots := make([]types.Dot, 0, w*h)

	for x := uint8(0); x <= s.zeroedBottomRight.X; x++ {
		for y := uint8(0); y <= s.zeroedBottomRight.Y; y++ {
			dots = append(dots, s.Absolute(x, y))
		}
	}

	return dots
}
