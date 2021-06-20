package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Dot_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input  []byte
		expect Dot
		err    bool
	}{
		{
			input:  []byte("1"),
			expect: Dot{},
			err:    true,
		},
		{
			input:  []byte("11"),
			expect: Dot{},
			err:    true,
		},
		{
			input: []byte("[2,3]"),
			expect: Dot{
				X: 2,
				Y: 3,
			},
			err: false,
		},
		{
			input: []byte("[255,255]"),
			expect: Dot{
				X: 255,
				Y: 255,
			},
			err: false,
		},
		{
			input:  nil,
			expect: Dot{},
			err:    true,
		},
		{
			input:  []byte("sadsaasdsasad"),
			expect: Dot{},
			err:    true,
		},
		{
			input:  []byte("[23]"),
			expect: Dot{},
			err:    true,
		},
		{
			input:  []byte("[23,32,54]"),
			expect: Dot{},
			err:    true,
		},
		{
			input:  []byte("[23, x]"),
			expect: Dot{},
			err:    true,
		},
		{
			input:  []byte("[y, x]"),
			expect: Dot{},
			err:    true,
		},
		{
			input:  []byte(`[123,"y"]`),
			expect: Dot{},
			err:    true,
		},
		{
			input:  []byte(`["y","x"]`),
			expect: Dot{},
			err:    true,
		},
		{
			input:  []byte(`"[y,x]"`),
			expect: Dot{},
			err:    true,
		},
		{
			input:  []byte("[y,x]"),
			expect: Dot{},
			err:    true,
		},
	}

	for i, test := range tests {
		var d Dot
		err := json.Unmarshal(test.input, &d)
		if test.err {
			require.NotNil(t, err, "Failed test number %d", i+1)
		} else {
			require.Nil(t, err, "Failed test number %d", i+1)
		}
		require.Equal(t, test.expect, d, "Failed test number %d", i+1)
	}
}

func Test_ObjectType_UnmarshalJSON(t *testing.T) {
	const msgFormat = "Failed test number %d"

	tests := []struct {
		input    []byte
		expected ObjectType
		err      bool
	}{
		{
			input:    []byte(`"snake"`),
			expected: ObjectTypeSnake,
			err:      false,
		},
		{
			input:    []byte(`"watermelon"`),
			expected: ObjectTypeWatermelon,
			err:      false,
		},
		{
			input:    []byte(`"corpse"`),
			expected: ObjectTypeCorpse,
			err:      false,
		},
		{
			input:    []byte(`"wall"`),
			expected: ObjectTypeWall,
			err:      false,
		},
		{
			input:    []byte(`"mouse"`),
			expected: ObjectTypeMouse,
			err:      false,
		},
		{
			input:    []byte(`"apple"`),
			expected: ObjectTypeApple,
			err:      false,
		},
		{
			input: []byte(`invalid value`),
			err:   true,
		},
		{
			input:    []byte(`"unknown"`),
			expected: ObjectTypeUnknown,
			err:      false,
		},
	}

	for i, test := range tests {
		var actualObjectType ObjectType
		err := json.Unmarshal(test.input, &actualObjectType)
		if test.err {
			require.NotNil(t, err, msgFormat, i+1)
		} else {
			require.Nil(t, err, msgFormat, i+1)
		}
		require.Equal(t, test.expected, actualObjectType, msgFormat, i+1)
	}
}
