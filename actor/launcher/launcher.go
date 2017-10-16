package launcher

import (
	"log"
	reflect "reflect"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/remote"
)

var (
	globalRegs = map[string]*actor.Props{}
	localRegs  = map[string]*actor.Props{}
)

func RegisterLocal(name string, prop *actor.Props) {
	if _, ok := localRegs[name]; ok {
		log.Fatalln(name, "exist")
	}
	localRegs[name] = prop
}

func RegisterGlobal(name string, prop *actor.Props) {
	if _, ok := globalRegs[name]; ok {
		log.Fatalln(name, "exist")
	}
	globalRegs[name] = prop
}

func AllGlobalNames() []string {
	ret := []string{}
	for name := range globalRegs {
		ret = append(ret, name)
	}
	return ret
}
func AllLocalNames() []string {
	ret := []string{}
	for name := range localRegs {
		ret = append(ret, name)
	}
	return ret
}
func RegToCluster(names []string) {
	for _, name := range names {
		if prop, ok := globalRegs[name]; ok {
			remote.Register(name, prop)
		}
	}
}

func StartLocals(names []string) {
	for _, name := range names {
		if prop, ok := localRegs[name]; ok {
			actor.SpawnNamed(prop, name)
		}
	}
}

func LogContext(ctx actor.Context) {
	log.Println(ctx.Self(), reflect.TypeOf(ctx.Message()), ctx.Message(), ctx.Sender())
}

//func RegSelf(pid *actor.PID) {
//	termPID, e := cluster.Get("term", "term")
//	if e != nil {
//		log.Fatal(e)
//	}
//	termPID.Tell(&term.RegLocal{})
//}
