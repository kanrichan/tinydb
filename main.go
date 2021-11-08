package main

import "fmt"

func main() {
	sto, err := NewCachingMiddleware(JSONStorage("test.json"))
	if err != nil {
		panic(err)
	}
	conn := TinyDB(sto).SetTable("xxx")
	conn.Exec(Insert(DocumentType{"123": 123}))
	fmt.Println(conn.Exec(Search().All()))
	fmt.Println(conn.Exec(Search().WhereEqual("123", 456)))
	conn.Exec(Remove().WhereEqual("123", 456))
	conn.Exec(Update(DocumentType{"234": 233}).All())
}
