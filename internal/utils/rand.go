package utils

import (
	"math/rand"
	"sync"
)

// NewRand returns a new Rand that uses a locked rand.Source.
//
// The Go rand.Source interface is not safe for concurrent use:
//
// - Docs: https://golang.org/pkg/math/rand/#Source
// - Discuss: https://github.com/golang/go/issues/3611
// - Solution: https://github.com/mesosphere/mesos-dns/pull/317
func NewRand(clock Clock) *rand.Rand {
	return rand.New(&lockedSource{
		src: rand.NewSource(clock.Now().UnixNano()),
	})
}

// lockedSource wraps a rand.Source with a sync.Mutex for synchronization.
type lockedSource struct {
	src rand.Source
	mux sync.Mutex
}

func (r *lockedSource) Int63() (n int64) {
	r.mux.Lock()
	n = r.src.Int63()
	r.mux.Unlock()
	return
}

func (r *lockedSource) Seed(seed int64) {
	r.mux.Lock()
	r.src.Seed(seed)
	r.mux.Unlock()
}
