package mtimer

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
