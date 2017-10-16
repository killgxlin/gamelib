package util

import "math/rand"

func Max(a ...int) int {
	if len(a) == 0 {
		panic("empty slice")
	}
	max := a[0]
	for i, n := range a {
		if i == 0 {
			continue
		}
		if n > max {
			max = n
		}
	}
	return max
}
func Min(a ...int) int {
	if len(a) == 0 {
		panic("empty slice")
	}
	min := a[0]
	for i, n := range a {
		if i == 0 {
			continue
		}
		if n < min {
			min = n
		}
	}
	return min
}
func Cnm(m, n int) []int {
	indexes := []int{}
	for i := 0; i < n; i++ {
		indexes = append(indexes, i)
	}
	rets := []int{}
	for i := 0; i < m && len(indexes) > 0; i++ {
		r := rand.Intn(len(indexes))
		selected := indexes[r]
		rets = append(rets, selected)
		indexes = append(indexes[0:r], indexes[r+1:]...)
		//indexes = EraseInt(selected, indexes)

	}
	return rets
}
