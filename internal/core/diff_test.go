package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_diff(t *testing.T) {
	const errFormat = "test case %d"

	tests := []struct {
		have, want, diff map[int]int
	}{
		{
			have: map[int]int{
				1: 1,
				2: 1,
				3: 1,
				7: 1,
				9: 5,
			},
			want: map[int]int{
				2: 4,
				8: 5,
				9: 3,
			},
			diff: map[int]int{
				1: -1,
				2: 3,
				3: -1,
				7: -1,
				8: 5,
				9: -2,
			},
		},
		{
			have: map[int]int{
				1: 1,
				2: 5,
				3: 1,
			},
			want: map[int]int{
				4: 4,
				5: 5,
				6: 3,
			},
			diff: map[int]int{
				1: -1,
				2: -5,
				3: -1,
				4: 4,
				5: 5,
				6: 3,
			},
		},
		{
			have: map[int]int{
				1: 3,
			},
			want: map[int]int{
				1: 3,
			},
			diff: map[int]int{},
		},
		{
			have: map[int]int{
				1: 10,
				2: 8,
			},
			want: map[int]int{
				1: 10,
				2: 12,
			},
			diff: map[int]int{
				2: 4,
			},
		},
	}

	for i, test := range tests {
		actualDiff := diff(test.have, test.want)
		assert.Equal(t, test.diff, actualDiff, errFormat, i+1)
	}
}

func Test_invertDiff(t *testing.T) {
	const errFormat = "test case %d"

	tests := []struct {
		diff, inverted map[int]int
	}{
		{
			diff: map[int]int{
				1: -1,
				2: 3,
				3: -1,
				7: -1,
				8: 5,
				9: -2,
			},
			inverted: map[int]int{
				1: 1,
				2: -3,
				3: 1,
				7: 1,
				8: -5,
				9: 2,
			},
		},
		{
			diff: map[int]int{
				1: -1,
				2: -5,
				3: -1,
				4: 4,
				5: 5,
				6: 3,
			},
			inverted: map[int]int{
				1: 1,
				2: 5,
				3: 1,
				4: -4,
				5: -5,
				6: -3,
			},
		},
		{
			diff: map[int]int{
				1: 10,
				2: 8,
			},
			inverted: map[int]int{
				1: -10,
				2: -8,
			},
		},
	}

	for i, test := range tests {
		actualInverted := invertDiff(test.diff)
		assert.Equal(t, test.inverted, actualInverted, errFormat, i+1)
	}
}

func Test_stateBotsNumber(t *testing.T) {
	const errFormat = "test case %d"

	tests := []struct {
		state  map[int]int
		number int
	}{
		{
			state: map[int]int{
				1: 1,
				2: 1,
				3: 1,
				7: 1,
				9: 5,
			},
			number: 9,
		},
		{
			state: map[int]int{
				1: 1,
				2: 5,
				3: 1,
			},
			number: 7,
		},
		{
			state: map[int]int{
				1: 3,
			},
			number: 3,
		},
		{
			state: map[int]int{
				1: 10,
				2: 8,
			},
			number: 18,
		},
	}

	for i, test := range tests {
		actualNumber := stateBotsNumber(test.state)
		assert.Equal(t, test.number, actualNumber, errFormat, i+1)
	}
}

func Test_diffStats(t *testing.T) {
	const errFormat = "test case %d"

	tests := []struct {
		diff   map[int]int
		add    int
		remove int
	}{
		{
			diff: map[int]int{
				1: -1,
				2: 3,
				3: -1,
				7: -1,
				8: 5,
				9: -2,
			},
			add:    8,
			remove: 5,
		},
		{
			diff: map[int]int{
				1: -1,
				2: -5,
				3: -1,
				4: 4,
				5: 5,
				6: 3,
			},
			add:    12,
			remove: 7,
		},
		{
			diff: map[int]int{
				1: 10,
				2: 8,
			},
			add:    18,
			remove: 0,
		},
	}

	for i, test := range tests {
		actualAdd, actualRemove := diffStats(test.diff)
		assert.Equal(t, test.add, actualAdd, errFormat, i+1)
		assert.Equal(t, test.remove, actualRemove, errFormat, i+1)
	}
}
