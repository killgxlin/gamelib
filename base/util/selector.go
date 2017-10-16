package util

import "path"

func SelectArr(channels []string, pattern string) map[string]bool {
	ret := map[string]bool{}
	for _, chn := range channels {
		if ok, _ := path.Match(pattern, chn); ok {
			ret[chn] = true
		}
	}
	return ret
}
func SelectMap(channels map[string]bool, pattern string) map[string]bool {
	ret := map[string]bool{}
	for chn := range channels {
		if ok, _ := path.Match(pattern, chn); ok {
			ret[chn] = true
		}
	}
	return ret
}
