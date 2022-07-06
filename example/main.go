package main

import (
	"fmt"

	tiny "github.com/Yiwen-Chan/tinydb"
)

type student struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func main() {
	storage, err := tiny.JSONStorage("test.json")
	if err != nil {
		panic(err)
	}
	database, err := tiny.TinyDB(storage)
	if err != nil {
		panic(err)
	}
	table := tiny.GetTable[student](database, "student")
	err = table.Insert(student{001, "test"})
	if err != nil {
		panic(err)
	}
	err = table.Update(func(s student) student { s.ID = 002; return s }, func(s student) bool { return true })
	if err != nil {
		panic(err)
	}
	stu, err := table.Select(func(s student) bool { return true })
	if err != nil {
		panic(err)
	}
	fmt.Println(stu)
	_, err = table.Delete(func(s student) bool { return true })
	if err != nil {
		panic(err)
	}
}
