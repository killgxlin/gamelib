package protobuf

import (
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/gogo/protobuf/proto"
)

type Coder struct {
	idType map[int32]reflect.Type // reflect.Type is *struct
	typeId map[reflect.Type]int32 // reflect.Type is *struct
}

func NewHandler() *Coder {
	h := &Coder{
		idType: make(map[int32]reflect.Type),
		typeId: make(map[reflect.Type]int32),
	}

	return h
}

func (h *Coder) RegTypes(pname string, valueId map[string]int32) {
	for sid, nid := range valueId {
		sname := strings.Replace(sid, "ID_", pname+".", 1)
		mT := proto.MessageType(sname)
		if mT == nil {
			log.Panicln("type not exist ", sname)
		}
		h.idType[nid] = mT
		h.typeId[mT] = nid
	}
}
func (h *Coder) Encode(m proto.Message) (buf []byte, e error) {
	mT := reflect.ValueOf(m).Type()
	mid, ok := h.typeId[mT]
	if !ok {
		log.Panicln("invalid msg type", mT, h.typeId)
		//e = fmt.Errorf("Handle.Encode:invalid msg type mT:%v, typeId:%v", mT, h.typeId)
		return
	}
	payload, e := proto.Marshal(m)
	if e != nil {
		return
	}
	buf = make([]byte, 4, 4+len(payload))
	binary.LittleEndian.PutUint32(buf, uint32(mid))
	buf = append(buf, payload...)
	return
}

func (h *Coder) Decode(buf []byte) (m proto.Message, e error) {
	if len(buf) < 4 {
		e = fmt.Errorf("len little then 4", buf)
		return
	}
	mid := int32(binary.LittleEndian.Uint32(buf[:4]))
	mT, ok := h.idType[mid]
	if !ok {
		e = fmt.Errorf("Handle.Decode:invalid msg id mid:%v, typeId:%v", mid, h.idType)
		return
		//log.Panicln("invalid msg id", mid, h.idType)
	}
	mV := reflect.New(mT.Elem())
	m = mV.Interface().(proto.Message)
	e = proto.Unmarshal(buf[4:], m)
	if e != nil {
		e = fmt.Errorf("Handle.Decode: err when unmarshal mid:%v, err:%v", mid, e)
	}
	return
}
