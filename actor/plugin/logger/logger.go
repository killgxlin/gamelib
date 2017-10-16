package logger

import (
	"log"
	"reflect"

	"github.com/AsynkronIT/protoactor-go/actor"
)

var (
	enableLog = true
)

func Enable(enable bool) {
	enableLog = enable
}

func MsgLogger(next actor.ActorFunc) actor.ActorFunc {
	return func(ctx actor.Context) {
		if enableLog {
			log.Printf("self:%v(%v) msg:%v(%v) sender:%v(%v)",
				ctx.Self(), reflect.TypeOf(ctx.Self()),
				ctx.Message(), reflect.TypeOf(ctx.Message()),
				ctx.Sender(), reflect.TypeOf(ctx.Sender()),
			)
		}
		next(ctx)
	}
}
