package engine

import (
	"sync/atomic"
	"unsafe"

	"github.com/ivan1993spb/snake-bot/internal/types"
)

type HashmapSight struct {
	sight  Sight
	values []unsafe.Pointer
}

func NewHashmapSight(s Sight) *HashmapSight {
	w := int(s.zeroedBottomRight.X) + 1
	h := int(s.zeroedBottomRight.Y) + 1
	values := make([]unsafe.Pointer, w*h)

	return &HashmapSight{
		sight:  s,
		values: values,
	}
}

func (h *HashmapSight) hash(dot types.Dot) int {
	if !h.sight.Seen(dot) {
		return -1
	}

	x, y := h.sight.Relative(dot)
	hash := int(y)*(int(h.sight.zeroedBottomRight.X)+1) + int(x)

	return hash
}

func (h *HashmapSight) unhash(i int) types.Dot {
	if i < 0 {
		return types.Dot{}
	}

	x := uint8(i % (int(h.sight.zeroedBottomRight.X) + 1))
	y := uint8(i / (int(h.sight.zeroedBottomRight.X) + 1))

	return h.sight.Absolute(x, y)
}

func (h *HashmapSight) Access(dot types.Dot) (interface{}, bool) {
	if hash := h.hash(dot); hash > -1 {
		p := atomic.LoadPointer(&h.values[hash])
		if uintptr(p) == uintptr(0) {
			return nil, false
		}
		return *(*interface{})(p), true
	}
	return nil, false
}

func (h *HashmapSight) AccessDefault(dot types.Dot,
	defaultValue interface{}) interface{} {
	if v, ok := h.Access(dot); ok {
		return v
	}
	return defaultValue

}

func (h *HashmapSight) Assign(dot types.Dot, v interface{}) {
	if hash := h.hash(dot); hash > -1 {
		atomic.SwapPointer(&h.values[hash], unsafe.Pointer(&v))
	}
}

func (h *HashmapSight) Flush() {
	for i := 0; i < len(h.values); i++ {
		atomic.SwapPointer(&h.values[i], unsafe.Pointer(uintptr(0)))
	}
}

func (h *HashmapSight) ForEach(f func(dot types.Dot, v interface{})) {
	for i := 0; i < len(h.values); i++ {
		p := atomic.LoadPointer(&h.values[i])
		if uintptr(p) != uintptr(0) {
			f(h.unhash(i), *(*interface{})(p))
		}
	}
}

func (h *HashmapSight) Reflect() *HashmapSight {
	return &HashmapSight{
		sight:  h.sight,
		values: make([]unsafe.Pointer, len(h.values)),
	}
}
