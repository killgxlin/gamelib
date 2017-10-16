package string

import (
	"bufio"
	"io"
	"gamelib/base/net"
)

type MsgReadWriter struct {
}

func NewReadWriter() *MsgReadWriter {
	return &MsgReadWriter{}
}

func (rw *MsgReadWriter) ReadMsgWithLimit(r io.Reader, limit int) (net.Message, error) {
	br := bufio.NewReader(r)
	m, _, e := br.ReadLine()
	if e != nil {
		m = nil
	}
	return string(m), e
}
func (rw *MsgReadWriter) WriteMsg(w io.Writer, m net.Message) error {
	bw := bufio.NewWriter(w)
	_, e := bw.WriteString(m.(string) + "\n")
	return e
}
