package engine

import (
	"fmt"

	"github.com/ivan1993spb/snake-bot/internal/types"
)

type Area struct {
	Width  uint8
	Height uint8
}

func NewArea(width, height uint8) Area {
	return Area{
		Width:  width,
		Height: height,
	}
}

func (a Area) Navigate(dot types.Dot) []types.Dot {
	dots := make([]types.Dot, 0, 4)

	// North
	if dot.Y > 0 {
		dots = append(dots, types.Dot{
			X: dot.X,
			Y: dot.Y - 1,
		})
	} else {
		dots = append(dots, types.Dot{
			X: dot.X,
			Y: a.Height - 1,
		})
	}

	// South
	dots = append(dots, types.Dot{
		X: dot.X,
		Y: (dot.Y + 1) % a.Height,
	})

	// East
	dots = append(dots, types.Dot{
		X: (dot.X + 1) % a.Width,
		Y: dot.Y,
	})

	// West
	if dot.X > 0 {
		dots = append(dots, types.Dot{
			X: dot.X - 1,
			Y: dot.Y,
		})
	} else {
		dots = append(dots, types.Dot{
			X: a.Width - 1,
			Y: dot.Y,
		})
	}

	return dots
}

func (a Area) FindDirection(from, to types.Dot) types.Direction {
	if from == to || !a.Fits(from) || !a.Fits(to) {
		return types.DirectionNorth
	}

	var (
		minDiffX, minDiffY uint8
		dirX, dirY         types.Direction
	)

	if from.X > to.X {
		minDiffX = from.X - to.X
		dirX = types.DirectionWest

		if diff := a.Width - from.X + to.X; diff < minDiffX {
			minDiffX = diff
			dirX = types.DirectionEast
		}
	} else {
		minDiffX = to.X - from.X
		dirX = types.DirectionEast

		if diff := a.Width - to.X + from.X; diff < minDiffX {
			minDiffX = diff
			dirX = types.DirectionWest
		}
	}

	if from.Y > to.Y {
		minDiffY = from.Y - to.Y
		dirY = types.DirectionNorth

		if diff := a.Height - from.Y + to.Y; diff < minDiffY {
			minDiffY = diff
			dirY = types.DirectionSouth
		}
	} else {
		minDiffY = to.Y - from.Y
		dirY = types.DirectionSouth

		if diff := a.Height - to.Y + from.Y; diff < minDiffY {
			minDiffY = diff
			dirY = types.DirectionNorth
		}
	}

	if minDiffX > minDiffY {
		return dirX
	}

	return dirY
}

func (a Area) FitDistance(distance, divisor, gap uint8) uint8 {
	if distance >= a.Width/divisor {
		distance = a.Width / divisor
		if gap > 0 && a.Width%divisor == 0 {
			distance -= gap
		}
	}

	if distance >= a.Height/divisor {
		distance = a.Height / divisor
		if gap > 0 && a.Height%divisor == 0 {
			distance -= gap
		}
	}

	return distance
}

func (a Area) Fits(dot types.Dot) bool {
	return a.Width > dot.X && a.Height > dot.Y
}

func (a Area) String() string {
	return fmt.Sprintf("Area[%d, %d]", a.Width, a.Height)
}
