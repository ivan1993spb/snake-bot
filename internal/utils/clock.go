package utils

import "time"

//counterfeiter:generate . Clock
type Clock interface {
	After(d time.Duration) <-chan time.Time
	Now() time.Time
}

var RealClock Clock = realClock{}

type realClock struct{}

func (realClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

func (realClock) Now() time.Time {
	return time.Now()
}

var NeverClock Clock = neverClock{}

type neverClock struct{}

func (neverClock) After(d time.Duration) <-chan time.Time {
	return make(chan time.Time)
}

func (neverClock) Now() time.Time {
	return time.Time{}
}

var ImmediatelyClock Clock = immediatelyClock{}

type immediatelyClock struct{}

func (immediatelyClock) After(d time.Duration) <-chan time.Time {
	ch := make(chan time.Time, 1)
	ch <- time.Time{}
	close(ch)
	return ch
}

func (immediatelyClock) Now() time.Time {
	return time.Now()
}
