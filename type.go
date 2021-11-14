package tinydb

import (
	"os"
	"sync"
)

// Storage An interface of Storage & Middleware.
// Should implement the method of Read() Write() Close().
type Storage interface {
	Read() (TinyTabs, error)
	Write(TinyTabs) error
	Close() error
}

// database the TinyDB class.
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

// TinyTabs
type TinyTabs map[string]TinyDocs

// TinyDocs
type TinyDocs map[int]TinyRecs

// TinyRecs
type TinyRecs map[string]TinyData

// TinyData
type TinyData interface{}

// Selector The function returns the selected serial number.
type Selector func(docs TinyDocs) []int

type TinyRecsIter struct {
	index int
	item  []TinyRecs
}

// TinyRecsArr The arrray of TinyRecs.
type TinyRecsArr []TinyRecs
