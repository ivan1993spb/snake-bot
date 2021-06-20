package engine

import (
	"sync/atomic"
	"unsafe"

	"github.com/ivan1993spb/snake-bot/internal/types"
)

type Map struct {
	fields [][]unsafe.Pointer
	area   Area
}

func NewMap(a Area) *Map {
	m := make([][]unsafe.Pointer, a.Height)

	for y := uint8(0); y < a.Height; y++ {
		m[y] = make([]unsafe.Pointer, a.Width)

		for x := uint8(0); x < a.Width; x++ {
			m[y][x] = unsafe.Pointer(uintptr(0))
		}
	}

	return &Map{
		fields: m,
		area:   a,
	}
}

func (m *Map) LookAround(s Sight) *HashmapSight {
	if m.area != s.area {
		return nil
	}
	objects := NewHashmapSight(s)
	for _, dot := range s.Dots() {
		if object, ok := m.GetObject(dot); ok {
			objects.Assign(dot, object)
		}
	}
	return objects
}

func (m *Map) GetObject(dot types.Dot) (*types.Object, bool) {
	if !m.area.Fits(dot) {
		return nil, false
	}

	p := atomic.LoadPointer(&m.fields[dot.Y][dot.X])

	if fieldIsEmpty(p) {
		return nil, false
	}

	container := (*types.Object)(p)

	return container, true
}

// fieldIsEmpty returns true if the pointer p is empty
func fieldIsEmpty(p unsafe.Pointer) bool {
	return uintptr(p) == uintptr(0)
}

func (m *Map) SaveObject(object *types.Object) {
	for _, dot := range object.GetDots() {
		if m.area.Fits(dot) {
			atomic.SwapPointer(&m.fields[dot.Y][dot.X], unsafe.Pointer(object))
		}
	}
}

func (m *Map) Clear(dots []types.Dot) {
	for _, dot := range dots {
		if m.area.Fits(dot) {
			atomic.SwapPointer(
				&m.fields[dot.Y][dot.X],
				unsafe.Pointer(uintptr(0)),
			)
		}
	}
}
