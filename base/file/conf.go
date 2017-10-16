package file

import (
	"encoding/json"
	"io/ioutil"
	"gamelib/base/util"
	"os"
)

func ReadFile(path string) []byte {
	pwd, e := os.Getwd()
	util.PanicOnErr(e)

	b, e := ioutil.ReadFile(pwd + path)
	util.PanicOnErr(e)
	return b
}

func ReadJson(path string, obj interface{}) {
	b := ReadFile(path)
	e := json.Unmarshal(b, obj)
	util.PanicOnErr(e)
}
