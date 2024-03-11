package core_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ivan1993spb/snake-bot/internal/core"
	"github.com/ivan1993spb/snake-bot/internal/core/corefakes"
	"github.com/ivan1993spb/snake-bot/internal/utils"
)

func Test_Core(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	botsLimit := 51
	factory := &corefakes.FakeBotOperatorFactory{}
	factory.NewReturns(&corefakes.FakeBotOperator{})

	params := &core.Params{
		BotsLimit:          botsLimit,
		BotOperatorFactory: factory,
		Clock:              utils.NeverClock,
	}

	c := core.NewCore(params)
	c.Run(ctx)

	// initial state
	state := map[int]int{
		1: 5,
		2: 4,
		3: 3,
		4: 2,
		5: 1,
	}

	callCount := 0

	t.Run("apply initial state", func(t *testing.T) {
		actual, err := c.SetState(ctx, state)
		require.NoError(t, err)
		require.Equal(t, state, actual)
		require.Equal(t, state, c.GetState(ctx))
		callCount += 15
		require.Equal(t, callCount, factory.NewCallCount())
	})

	t.Run("change first one", func(t *testing.T) {
		state[1] = 10
		actual, err := c.SetOne(ctx, 1, 10)
		require.NoError(t, err)
		require.Equal(t, state, actual)
		require.Equal(t, state, c.GetState(ctx))
		callCount += 5
		require.Equal(t, callCount, factory.NewCallCount())
	})

	t.Run("apply state, add and remove", func(t *testing.T) {
		delete(state, 1)
		delete(state, 2)
		state[3] = 3
		delete(state, 4)
		state[5] = 3
		state[6] = 3

		actual, err := c.SetState(ctx, map[int]int{
			1: 0,
			2: 0,
			3: 3,
			4: 0,
			5: 3,
			6: 3,
		})
		require.NoError(t, err)
		require.Equal(t, state, actual)
		require.Equal(t, state, c.GetState(ctx))
		callCount += 5
		require.Equal(t, callCount, factory.NewCallCount())
	})

	t.Run("apply state exceed limit", func(t *testing.T) {
		actual, err := c.SetState(ctx, map[int]int{
			7: botsLimit,
			2: 1,
		})
		require.Error(t, err)
		require.ErrorIs(t, err, core.ErrRequestedTooManyBots)
		require.Nil(t, actual)
		require.Equal(t, state, c.GetState(ctx))
		// callCount not changed
		require.Equal(t, callCount, factory.NewCallCount())
	})

	t.Run("change one exceed limit", func(t *testing.T) {
		actual, err := c.SetOne(ctx, 7, botsLimit)
		require.Error(t, err)
		require.ErrorIs(t, err, core.ErrRequestedTooManyBots)
		require.Nil(t, actual)
		require.Equal(t, state, c.GetState(ctx))
		// callCount not changed
		require.Equal(t, callCount, factory.NewCallCount())
	})
}
