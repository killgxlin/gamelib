package util

import (
	"os"
	"path/filepath"
	"strings"
)

func GetAbsPath(path string) string {
	dir, e := filepath.Abs(filepath.Dir(os.Args[0]) + path)
	PanicOnErr(e)
	return strings.Replace(dir, "\\", "/", -1)
}
