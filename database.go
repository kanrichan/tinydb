package tinydb

import (
	"io"
	"sync"
)

// Database the TinyDB database.
type Database struct {
	sync.Mutex
	table   string
	storage Storage
}

// TinyDB create a new database with the interface of storage.
func TinyDB(storage Storage) (*Database, error) {
	var database = &Database{sync.Mutex{}, "_default", storage}
	database.Lock()
	defer database.Unlock()
	var data = make(map[string][]interface{})
	err := storage.Read(&data)
	if err == io.EOF {
		err := storage.Write(map[string][]interface{}{"_default": nil})
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return database, nil
}

// Close the database.
func (db *Database) Close() error {
	return db.storage.Close()
}

// Tables get the names of all tables in the database.
func (db *Database) Tables() ([]string, error) {
	var data = make(map[string][]interface{})
	err := db.storage.Read(&data)
	if err != nil {
		return nil, err
	}
	var out = make([]string, 0)
	for k := range data {
		out = append(out, k)
	}
	return out, nil
}
