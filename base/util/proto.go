package util

import (
	"reflect"

	"github.com/gogo/protobuf/proto"
)

func Clone(m proto.Message) proto.Message {
	b, e := proto.Marshal(m)
	PanicOnErr(e)
	m1 := reflect.New(reflect.ValueOf(m).Type().Elem()).Interface().(proto.Message)
	e = proto.Unmarshal(b, m1)
	PanicOnErr(e)

	return m1
}
