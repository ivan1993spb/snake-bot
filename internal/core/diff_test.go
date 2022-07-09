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

func Test_diffOne(t *testing.T) {
	const errFormat = "test case %d"

	tests := []struct {
		have map[int]int
		k    int
		v    int
		diff map[int]int
	}{
		{
			have: map[int]int{
				1: 10,
				2: 8,
				4: 32,
			},
			k: 12,
			v: 1,
			diff: map[int]int{
				12: 1,
			},
		},
		{
			have: map[int]int{
				1:  10,
				2:  8,
				4:  32,
				12: 21,
			},
			k: 12,
			v: 1,
			diff: map[int]int{
				12: -20,
			},
		},
		{
			have: map[int]int{
				1:  10,
				2:  8,
				4:  32,
				12: 21,
			},
			k: 2,
			v: 10,
			diff: map[int]int{
				2: 2,
			},
		},
	}

	for i, test := range tests {
		actualDiff := diffOne(test.k, test.v, test.have)
		assert.Equal(t, test.diff, actualDiff, errFormat, i+1)
	}
}
