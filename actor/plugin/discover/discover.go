package discover

import (
	mutil "gamelib/actor/util"
	"gamelib/base/util"
	"gamelib/etcd/discover"

	"github.com/AsynkronIT/protoactor-go/actor"
)

// middleward ----------------------------------------------------------------------------
type AddWatch struct {
	Lprefix string
	Typ     discover.WatchType
}

type WatchEvent struct {
	Lprefix string
	PIDs    []*actor.PID
}

// warning
func MakeWatcher() func(actor.ActorFunc) actor.ActorFunc {
	return func(next actor.ActorFunc) actor.ActorFunc {
		watchers := map[string]*discover.Watcher{}
		return func(ctx actor.Context) {
			switch m := ctx.Message().(type) {
			case *actor.Started:
			case *actor.Stopping, *actor.Restarting:
				for _, w := range watchers {
					w.Close()
				}
			case *AddWatch:
				discover.NewWatch("/services/"+m.Lprefix, m.Typ, func(prefix string, values []string) {
					we := &WatchEvent{Lprefix: m.Lprefix}
					//log.Println("watchcb", values)
					for _, v := range values {
						pid := mutil.KeyToPID2(v)
						we.PIDs = append(we.PIDs, pid)
					}
					ctx.Self().Tell(we)
				})
			}
			next(ctx)
		}
	}
}

// warning
func MakeStub(label string) func(actor.ActorFunc) actor.ActorFunc {
	return func(next actor.ActorFunc) actor.ActorFunc {
		var stub *discover.Stub
		return func(ctx actor.Context) {
			switch ctx.Message().(type) {
			case *actor.Started:
				s, e := discover.NewStub("/locks/"+label, "/services/"+label, mutil.PIDToKey2(ctx.Self()))
				util.PanicOnErr(e)
				stub = s
			case *actor.Stopping, *actor.Restarting:
				if stub != nil {
					stub.Close()
					stub = nil
				}
			}
			next(ctx)
		}
	}
}
