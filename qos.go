package qos

import (
	"github.com/rfyiamcool/backoff"
	"sync/atomic"
	"time"
)

type Executor struct {
	pending  int32
	MaxDelay time.Duration
	MinDelay time.Duration
}

func (e *Executor) Pending() int32 {
	return atomic.LoadInt32(&e.pending)
}

func (e *Executor) Execute(f func() error) {
	atomic.AddInt32(&e.pending, 1)
	go func() {
		defer atomic.AddInt32(&e.pending, -1)

		b := backoff.NewBackOff(
			backoff.WithMinDelay(e.MinDelay),
			backoff.WithMaxDelay(e.MaxDelay),
			backoff.WithFactor(2),
		)

		for {
			if err := f(); err == nil {
				return
			}

			b.Sleep()
		}
	}()
}
