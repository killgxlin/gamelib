package cmdhandler

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"

	gsw "github.com/mattn/go-shellwords"
)

type Handler func(args []string, fs *flag.FlagSet, out io.Writer, ctx interface{})

type CmdHandler struct {
	handles map[string]Handler
}

func (ch *CmdHandler) AllCommand() []string {
	var cmds []string
	for cmd := range ch.handles {
		cmds = append(cmds, cmd)
	}
	return cmds
}

func (ch *CmdHandler) Register(typ string, h Handler) {
	if _, ok := ch.handles[typ]; ok {
		log.Panic(typ, "exist")
	}
	ch.handles[typ] = h
}

func (ch *CmdHandler) Handle(cmd string, ctx interface{}) (ret string, e error) {
	defer func() {
		if r := recover(); r != nil {
			e = fmt.Errorf("%v", r)
		}
	}()
	args, e := gsw.Parse(cmd)
	if e != nil {
		return "", e
	}

	h, ok := ch.handles[args[0]]
	if !ok {
		return "", fmt.Errorf("handler of %v no exist", args[0])
	}

	b := &bytes.Buffer{}
	fs := flag.NewFlagSet(args[0], flag.ContinueOnError)
	fs.SetOutput(b)

	h(args[1:], fs, b, ctx)

	if b.Len() > 0 {
		b.Truncate(b.Len() - 1)
	}
	ret = b.String()
	e = nil

	return
}

func NewCmdHandler() *CmdHandler {
	return &CmdHandler{
		handles: map[string]Handler{},
	}
}
