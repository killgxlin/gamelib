package profile

import (
	nutil "gamelib/base/net/util"
	"gamelib/base/util"
	"net/http"
	_ "net/http/pprof"
)

func StartWebTrace() {
	go func() {
		l, e := nutil.Listen("tcp", "localhost", 6060, 8080)
		util.PanicOnErr(e)
		http.Serve(l, nil)
	}()
}
