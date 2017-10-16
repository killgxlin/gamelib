package util

import (
	"bytes"
	"fmt"
)

func Concate(sep string, is ...interface{}) string {
	var b bytes.Buffer
	for i, a := range is {
		b.WriteString(fmt.Sprint(a))
		if i != len(is)-1 {
			b.WriteString(sep)
		}
	}
	return b.String()
}
