package tinydb

import (
	"os"
	"sync"
)

// Storage An interface of Storage & Middleware
// Should implement the method of Read() Write() Close()
type Storage interface {
	Read() (TinyTabs, error)
	Write(TinyTabs) error
	Close() error
}

// database the TinyDB class
type database struct {
	sync.Mutex
	table   string
	storage Storage
}

// StorageJSON Store the data in a JSON file.
type StorageJSON struct {
	mutex  sync.Mutex
	handle *os.File
}

// StorageMemory Store the data in a memory.
type StorageMemory struct {
	memory []byte
}

// MiddlewareCaching
type MiddlewareCaching struct {
	storage Storage
	cache   []byte
	count   int
	size    int
}

type TinyTabs map[string]TinyDocs
type TinyDocs map[int]TinyRecs
type TinyRecs map[string]TinyData
type TinyData interface{}

type Selector func(docs TinyDocs) []int

type TinyRecsIter struct {
	index int
	item  []TinyRecs
}

type TinyRecsArr []TinyRecs
