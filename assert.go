package main

import (
	"errors"
	"fmt"
)

type TinyValue struct {
	Num  float64
	Str  string
	Bool bool
	Map  map[string]*TinyValue
	Arr  []*TinyValue
}

func main() {
	a := Document{
		"123": map[string]interface{}{
			"123": map[string]interface{}{
				"123": 123.0,
			},
		},
	}
	fmt.Println(a.Get("123").Map["123"].Map["123"].Num)
}

func (doc Document) Get(key string) *TinyValue {
	return AssertJSONValue(doc[key])
}

func AssertJSONValue(v interface{}) *TinyValue {
	switch v.(type) {
	case float64:
		return &TinyValue{Num: v.(float64)}
	case string:
		return &TinyValue{Str: v.(string)}
	case bool:
		return &TinyValue{Bool: v.(bool)}
	case map[string]interface{}:
		return AssertJSONMap(v.(map[string]interface{}))
	case []interface{}:
		return AssertJSONArr(v.([]interface{}))
	case nil:
		return &TinyValue{}
	default:
		panic(errors.New("what's type???"))
	}
}

func AssertJSONMap(v map[string]interface{}) *TinyValue {
	var dict = map[string]*TinyValue{}
	for k := range v {
		dict[k] = AssertJSONValue(v[k])
	}
	return &TinyValue{Map: dict}
}

func AssertJSONArr(v []interface{}) *TinyValue {
	var arr = []*TinyValue{}
	for k := range v {
		arr[k] = AssertJSONValue(v[k])
	}
	return &TinyValue{Arr: arr}
}
