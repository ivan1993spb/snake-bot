package engine

import (
	"testing"

	"github.com/ivan1993spb/snake-bot/internal/types"
)

func Benchmark_HashmapSight_Assign(b *testing.B) {
	const (
		width  = 255
		height = 255

		distance = 100
	)

	a := NewArea(width, height)
	pos := types.Dot{
		X: 0,
		Y: 0,
	}
	s := NewSight(a, pos, distance)
	dots := s.Dots()
	l := len(dots)

	h := NewHashmapSight(s)

	val := struct {
		a, b, c int
		d, e, f interface{}
		g, h, i uint8
	}{}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dot := dots[i%l]
		h.Assign(dot, val)
	}
}

func Benchmark_Map_Assign(b *testing.B) {
	const (
		width  = 255
		height = 255

		distance = 100
	)

	a := NewArea(width, height)
	pos := types.Dot{
		X: 0,
		Y: 0,
	}
	s := NewSight(a, pos, distance)
	dots := s.Dots()
	l := len(dots)

	h := make(map[types.Dot]interface{}, l)

	val := struct {
		a, b, c int
		d, e, f interface{}
		g, h, i uint8
	}{}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dot := dots[i%l]
		h[dot] = val
	}
}

func Benchmark_HashmapSight_Access(b *testing.B) {
	const (
		width  = 255
		height = 255

		distance = 100
	)

	a := NewArea(width, height)
	pos := types.Dot{
		X: 0,
		Y: 0,
	}
	s := NewSight(a, pos, distance)
	dots := s.Dots()
	l := len(dots)

	h := NewHashmapSight(s)

	val := struct {
		a, b, c int
		d, e, f interface{}
		g, h, i uint8
	}{}

	for _, dot := range dots {
		h.Assign(dot, val)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dot := dots[i%l]
		_, _ = h.Access(dot)
	}
}

func Benchmark_Map_Access(b *testing.B) {
	const (
		width  = 255
		height = 255

		distance = 100
	)

	a := NewArea(width, height)
	pos := types.Dot{
		X: 0,
		Y: 0,
	}
	s := NewSight(a, pos, distance)
	dots := s.Dots()
	l := len(dots)

	h := make(map[types.Dot]interface{}, l)

	val := struct {
		a, b, c int
		d, e, f interface{}
		g, h, i uint8
	}{}

	for _, dot := range dots {
		h[dot] = val
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dot := dots[i%l]
		_, _ = h[dot]
	}
}
