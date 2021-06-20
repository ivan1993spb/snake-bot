package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-bot/internal/types"
)

func TestHashmapSight_hash(t *testing.T) {
	const msgFormat = "hash(%s) = %d"
	const (
		width  = 30
		height = 22
	)

	a := NewArea(width, height)
	d := uint8(4)
	for posX := uint8(0); posX < a.Width; posX++ {
		for posY := uint8(0); posY < a.Height; posY++ {
			pos := types.Dot{
				X: posX,
				Y: posY,
			}
			s := NewSight(a, pos, d)
			h := NewHashmapSight(s)

			for x := uint8(0); x < a.Width; x++ {
				for y := uint8(0); y < a.Height; y++ {
					dot := types.Dot{
						X: x,
						Y: y,
					}
					res := h.hash(dot)
					require.Truef(t, res == -1 || res < len(h.values),
						msgFormat, dot, res)
				}
			}
		}
	}
}

func TestHashmapSight_UniqHashes(t *testing.T) {
	const msgFormat = "hash(%s) = %d (the same %s)"

	const (
		width  = 255
		height = 255
	)

	a := NewArea(width, height)
	d := uint8(1)
	pos := types.Dot{
		X: 0,
		Y: 0,
	}
	s := NewSight(a, pos, d)
	h := NewHashmapSight(s)

	check := map[int]types.Dot{}

	for _, dot := range s.Dots() {
		hash := h.hash(dot)
		dot1, ok := check[hash]
		assert.Falsef(t, ok, msgFormat, dot, hash, dot1)
		check[hash] = dot
	}
}

func TestHashmapSight_AssignAccessFlush(t *testing.T) {
	const (
		width  = 30
		height = 22
	)

	a := NewArea(width, height)
	d := uint8(4)
	pos := types.Dot{
		X: 2,
		Y: 4,
	}
	pos1 := types.Dot{
		X: 3,
		Y: 4,
	}
	s := NewSight(a, pos, d)
	h := NewHashmapSight(s)
	h.Assign(pos, 2)

	val, ok := h.Access(pos)
	require.True(t, ok)
	require.Equal(t, 2, val)

	val1, ok := h.Access(pos1)
	require.False(t, ok)
	require.Nil(t, val1)

	h.Flush()

	val2, ok := h.Access(pos)
	require.False(t, ok)
	require.Nil(t, val2)
}

func TestHashmapSight_hash_unhash(t *testing.T) {
	const msgFormat = "hash(%s) = %d; unhash(%d) = %s"

	const (
		width  = 30
		height = 22
	)

	a := NewArea(width, height)
	d := uint8(4)
	pos := types.Dot{
		X: 1,
		Y: 0,
	}
	s := NewSight(a, pos, d)
	h := NewHashmapSight(s)

	for _, dot := range s.Dots() {
		hash := h.hash(dot)
		unhash := h.unhash(hash)
		require.Equalf(t, dot, unhash, msgFormat,
			dot, hash, hash, unhash)
	}
}

func TestHashmapSight_AccessReturnsFalseIfEmpty(t *testing.T) {
	const msgFormat = "case number %d"

	const (
		width  = 220
		height = 147
	)

	a := NewArea(width, height)
	d := uint8(50)
	pos := types.Dot{
		X: 1,
		Y: 0,
	}
	s := NewSight(a, pos, d)
	h := NewHashmapSight(s)

	for i, dot := range s.Dots() {
		_, ok := h.Access(dot)
		require.Falsef(t, ok, msgFormat, i+1)
		h.Assign(dot, "ok")
	}
}

func TestHashmapSight_CompareWithMap(t *testing.T) {
	const (
		width  = 220
		height = 147
	)

	a := NewArea(width, height)
	d := uint8(50)
	pos := types.Dot{
		X: 1,
		Y: 0,
	}
	s := NewSight(a, pos, d)
	h := NewHashmapSight(s)

	m := map[types.Dot]interface{}{}

	for i, dot := range s.Dots() {
		_, ok := h.Access(dot)
		assert.False(t, ok)

		m[dot] = i
		h.Assign(dot, i)
	}

	h.ForEach(func(dot types.Dot, v interface{}) {
		v1, ok := m[dot]
		assert.True(t, ok)
		assert.Equal(t, v, v1)
	})

	for dot, v := range m {
		v1, ok := h.Access(dot)
		assert.True(t, ok)
		assert.Equal(t, v, v1)
	}
}
