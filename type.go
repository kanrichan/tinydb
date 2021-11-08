package tinydb

import (
	"os"
	"sync"
)

// StorageData
type StorageData map[string]TableData

// TableData
type TableData map[int]Document

// Document
type Document map[string]interface{}

// Storage An interface of Storage & Middleware
// Should implement the method of Read() Write() Close()
type Storage interface {
	Read() (StorageData, error)
	Write(StorageData) error
	Close() error
}

// database the TinyDB class
type database struct {
	table string
	conn  chan *Request
	stop  chan bool
}

// StorageJSON Store the data in a JSON file.
type StorageJSON struct {
	mutex  sync.Mutex
	handle *os.File
}

// StorageMemory Store the data in a memory.
type StorageMemory struct {
	memory StorageData
}

// MiddlewareCaching
type MiddlewareCaching struct {
	storage Storage
	cache   StorageData
	count   int
	size    int
}

// Request
type Request struct {
	table     string
	operation func() func(storage Storage) ([]Document, error)
	condition func() func(doc Document) bool
	response  chan *Response
}

// Response
type Response struct {
	err  error
	docs []Document
}
