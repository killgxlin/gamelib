package typhandler

import (
	"fmt"
	"log"
	"reflect"
)

var (
	ErrNoHandle = fmt.Errorf("func no exist type")
)

type Handler struct {
	typeFunc map[string][]reflect.Value
}

func (h *Handler) RegByType(i interface{}, f interface{}) {
	iT := reflect.TypeOf(i)
	switch iT.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		iT = iT.Elem()
	}

	fT := reflect.TypeOf(f)

	if fT.NumIn() < 1 {
		log.Panicln("invalid arg type")
	}

	switch fT.NumOut() {
	case 0:
	case 1:
		eT := reflect.TypeOf((*error)(nil)).Elem()
		if !fT.Out(0).Implements(eT) {
			log.Panicln("invalid rpc ret type")
		}
	default:
		log.Panicln("invalid f type")
	}

	mT := fT.In(0)
	fVs, ok := h.typeFunc[iT.Name()]
	if ok {
		log.Panicln("func exist ", mT)
	}
	h.typeFunc[iT.Name()] = append(fVs, reflect.ValueOf(f))
}
func (h *Handler) RegByFunc(f interface{}) {
	fT := reflect.TypeOf(f)
	if fT.NumIn() < 1 {
		log.Panicln("invalid arg type")
	}

	aT := fT.In(0)
	switch aT.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		aT = aT.Elem()
	}

	fVs := h.typeFunc[aT.Name()]
	switch fT.NumOut() {
	case 0:
	case 1:
		eT := reflect.TypeOf((*error)(nil)).Elem()
		if !fT.Out(0).Implements(eT) {
			log.Panicln("invalid rpc ret type")
		}
	default:
		log.Panicln("invalid f type")
	}

	h.typeFunc[aT.Name()] = append(fVs, reflect.ValueOf(f))
}

func (h *Handler) Handle(m interface{}, arg ...interface{}) error {
	mV := reflect.ValueOf(m)

	aT := mV.Type()
	switch aT.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Ptr, reflect.Slice:
		aT = aT.Elem()
	}

	fV, ok := h.typeFunc[aT.Name()]
	if !ok {
		return ErrNoHandle
	}

	aV := []reflect.Value{mV}
	for _, a := range arg {
		aV = append(aV, reflect.ValueOf(a))
	}

	for _, fv := range fV {
		vV := fv.Call(aV)
		if len(vV) < 1 {
			continue
		}
		e, _ := vV[0].Interface().(error)
		if e != nil {
			return e
		}
	}
	return nil
}

func NewHandler() *Handler {
	return &Handler{
		typeFunc: map[string][]reflect.Value{},
	}
}
