package main

import (
	"fmt"

	tiny "github.com/Yiwen-Chan/tinydb"
)

func main() {
	db, err := tiny.TinyDB(tiny.CachingMiddleware(tiny.JSONStorage("test.json")))
	if err != nil {
		panic(err)
	}
	err = db.Insert(tiny.TinyRecs{"test": 123})
	if err != nil {
		panic(err)
	}
	iter, err := db.Search(tiny.All().OrderBy("test", true))
	if err != nil {
		panic(err)
	}
	for iter.HasNext() {
		fmt.Println(iter.Next().Get("test").Num)
	}
}
