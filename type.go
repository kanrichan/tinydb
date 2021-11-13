package main

import (
	"os"
	"sync"
)

// Storage An interface of Storage & Middleware
// Should implement the method of Read() Write() Close()
type Storage interface {
	Read() (TinyTabsMap, error)
	Write(TinyTabsMap) error
	Close() error
}

// database the TinyDB class
type database struct {
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
	memory TinyTabsMap
}

// MiddlewareCaching
type MiddlewareCaching struct {
	storage Storage
	cache   TinyTabsMap
	count   int
	size    int
}
