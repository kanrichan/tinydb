package tinydb

import (
	"io"
	"sync"
)

// database the TinyDB class.
type database struct {
	sync.Mutex
	table   string
	storage Storage
}

// TinyDB Create a new database with the interface of storage.
func TinyDB(storage Storage) (*database, error) {
	var database = &database{sync.Mutex{}, "_default", storage}
	database.Lock()
	defer database.Unlock()
	var data = make(map[string][]interface{})
	err := storage.Read(data)
	if err == io.EOF {
		storage.Write(map[string][]interface{}{"_default": nil})
	}
	return database, nil
}

func GetTable[T any](database *database, name string) *Table[T] {
	return &Table[T]{name, database}
}

func DropTable[T any](database *database, name string) error {
	database.Lock()
	defer database.Unlock()
	var data = make(map[string][]interface{})
	err := database.storage.Read(data)
	if err != nil {
		return err
	}
	delete(data, name)
	return database.storage.Write(data)
}
