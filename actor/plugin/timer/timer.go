package timer

import (
	"context"
	"sync"
	"time"

	"github.com/AsynkronIT/protoactor-go/actor"
)

type TimerEvent struct {
}

func MakeTimer(dur time.Duration) func(actor.ActorFunc) actor.ActorFunc {
	return func(next actor.ActorFunc) actor.ActorFunc {
		wg := &sync.WaitGroup{}
		con, can := context.WithCancel(context.TODO())
		return func(ctx actor.Context) {
			switch ctx.Message().(type) {
			case *actor.Started:
				wg.Add(1)
				go func() {
					t := time.NewTicker(dur)
					defer func() {
						t.Stop()
						wg.Done()
					}()
					for {
						select {
						case <-con.Done():
							return
						case <-t.C:
							ctx.Self().Tell(TimerEvent{})
						}
					}
				}()
			case *actor.Stopping, *actor.Restarting:
				can()
				wg.Wait()
			}
			next(ctx)
		}
	}
}

type Timer struct {
	context.Context
	context.CancelFunc
	*sync.WaitGroup
	time.Duration
}

func NewTimer(dur time.Duration) *Timer {
	tmr := &Timer{}
	tmr.Context, tmr.CancelFunc = context.WithCancel(context.TODO())
	tmr.WaitGroup = &sync.WaitGroup{}
	tmr.Duration = dur

	return tmr
}

func (t *Timer) OnStart(ctx actor.Context) {
	t.WaitGroup.Add(1)
	go func() {
		ticker := time.NewTicker(t.Duration)
		defer func() {
			ticker.Stop()
			t.WaitGroup.Done()
		}()
		for {
			select {
			case <-t.Context.Done():
				return
			case <-ticker.C:
				ctx.Self().Tell(TimerEvent{})
			}
		}
	}()
}
func (t *Timer) OnOtherMessage(ctx actor.Context, m interface{}) {
	switch m.(type) {
	case *actor.Stopping, *actor.Restarting:
		t.CancelFunc()
		t.WaitGroup.Wait()
	}
}
