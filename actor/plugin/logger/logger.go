package logger

import (
	"log"
	"reflect"

	"github.com/AsynkronIT/protoactor-go/actor"
)

var (
	enableLog = true
	filter    = map[reflect.Type]struct{}{}
)

func Enable(enable bool) {
	enableLog = enable
}

func Filter(t ...interface{}) {
	for _, t := range t {
		filter[reflect.TypeOf(t)] = struct{}{}
	}
}

func MsgLogger(next actor.ActorFunc) actor.ActorFunc {
	return func(ctx actor.Context) {
		if enableLog {
			mt := reflect.TypeOf(ctx.Message())
			if _, ok := filter[mt]; !ok {
				log.Printf("self:%v(%v) msg:%v(%v) sender:%v(%v)",
					ctx.Self(), reflect.TypeOf(ctx.Self()),
					ctx.Message(), mt,
					ctx.Sender(), reflect.TypeOf(ctx.Sender()),
				)
			}
		}
		next(ctx)
	}
}
