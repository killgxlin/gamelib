package timer

import (
	"context"
	"sync"
	"time"
)

type Timer struct {
	count  int
	dur    time.Duration
	ticker *time.Ticker
	f      func(cur time.Time, last bool)
	w      sync.WaitGroup
	ctx    context.Context
	can    context.CancelFunc
}

func NewTimer(dur time.Duration, count int, f func(cur time.Time, last bool)) *Timer {
	tmr := &Timer{
		count:  count,
		dur:    dur,
		f:      f,
		ticker: time.NewTicker(dur),
	}

	tmr.ctx, tmr.can = context.WithCancel(context.TODO())

	tmr.w.Add(1)
	go func() {
		defer func() {
			tmr.w.Done()
		}()
		for {
			select {
			case now := <-tmr.ticker.C:
				switch {
				case tmr.count > 0:
					tmr.count--
					f(now, tmr.count == 0)
					if tmr.count == 0 {
						tmr.ticker.Stop()
						tmr.can()
						return
					}
				case tmr.count == 0:
					f(now, false)
				default:
					panic("timer - value")
				}
			case <-tmr.ctx.Done():
				return
			}
		}
	}()

	return tmr
}

func (t *Timer) Stop() {
	if t == nil {
		return
	}

	t.ticker.Stop()
	t.can()
	t.w.Wait()
}
