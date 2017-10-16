package protobuf

import (
	"encoding/binary"
	"io"
	"gamelib/base/net"
	"gamelib/base/util"

	"github.com/gogo/protobuf/proto"
)

type MsgReadWriter struct {
	coder *Coder
}

func NewReadWriter() *MsgReadWriter {
	return &MsgReadWriter{coder: NewHandler()}
}

func (rw *MsgReadWriter) RegTypes(pname string, valueId map[string]int32) {
	rw.coder.RegTypes(pname, valueId)
}

func (rw *MsgReadWriter) ReadMsgWithLimit(r io.Reader, limit int) (m net.Message, e error) {
	sz := uint32(0)
	e = binary.Read(r, binary.BigEndian, &sz)
	if e != nil {
		return
	}
	if sz >= uint32(limit) {
		e = net.ErrSizeExceedLimit
		return
	}
	raw := make([]byte, sz)
	_, e = io.ReadFull(r, raw)
	if e != nil {
		return
	}

	return rw.coder.Decode(raw)
}
func (rw *MsgReadWriter) WriteMsg(w io.Writer, m net.Message) error {
	raw, e := rw.coder.Encode(m.(proto.Message))
	util.PanicOnErr(e)
	if e != nil {
		return e
	}
	sz := uint32(len(raw))

	e = binary.Write(w, binary.BigEndian, sz)
	if e != nil {
		return e
	}

	_, e = w.Write(raw)
	return e
}
