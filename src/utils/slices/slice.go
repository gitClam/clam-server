package slices

import (
	"reflect"
)

func ContainsInSlice(c interface{}, value interface{}) bool {
	m := ToMapSet(c)
	if m == nil {
		return false
	}
	return ContainsInMap(m, value)
}

func ContainsInMap(m map[interface{}]struct{}, key interface{}) bool {
	_, ok := m[key]
	return ok
}

func ToMapSet(i interface{}) map[interface{}]struct{} {
	// judge the validation of the input
	if i == nil {
		return nil
	}
	kind := reflect.TypeOf(i).Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil
	}

	// execute the convert
	v := reflect.ValueOf(i)
	m := make(map[interface{}]struct{}, v.Len())
	for j := 0; j < v.Len(); j++ {
		m[v.Index(j).Interface()] = struct{}{}
	}
	return m
}
