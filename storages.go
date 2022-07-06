package tinydb

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

// Storage An interface of Storage & Middleware.
// Should implement the method of Read() Write() Close().
type Storage interface {
	Read(any) error
	Write(any) error
	Close() error
}

// StorageJSON Store the data in a JSON file.
type StorageJSON struct {
	handle *os.File
}

// JSONStorage Create a new JSONStorage instance.
func JSONStorage(file string) (*StorageJSON, error) {
	var dir string
	i1 := strings.Index(file, `\`)
	i2 := strings.Index(file, `/`)
	if i1 != -1 || i2 != -1 {
		if i1 > i2 {
			dir = file[:i1]
		} else {
			dir = file[:i2]
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, err
			}
		}
	}
	fi, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	return &StorageJSON{handle: fi}, err
}

// Read Read data from JSON file.
func (sto *StorageJSON) Read(data any) error {
	sto.handle.Seek(0, 0)
	dec := json.NewDecoder(sto.handle)
	return dec.Decode(data)
}

// Write Write data to JSON file.
func (sto *StorageJSON) Write(data any) error {
	sto.handle.Truncate(0)
	sto.handle.Seek(0, 0)
	enc := json.NewEncoder(sto.handle)
	enc.SetIndent("", "    ")
	if data == nil {
		return errors.New("Nothing needs to be written")
	}
	return enc.Encode(data)
}

// Close Close the JSONStorage instance.
func (sto *StorageJSON) Close() error {
	return sto.handle.Close()
}

// StorageMemory Store the data in a memory.
type StorageMemory struct {
	memory []byte
}

// MemoryStorage Create a new MemoryStorage instance.
func MemoryStorage() (*StorageMemory, error) {
	return &StorageMemory{memory: []byte{}}, nil
}

// Read Read data from memory.
func (sto *StorageMemory) Read(data any) error {
	return json.Unmarshal(sto.memory, &data)
}

// Write Write data to memory.
func (sto *StorageMemory) Write(data any) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	sto.memory = b
	return nil
}

// Close Close the MemoryStorage instance.
func (sto *StorageMemory) Close() error {
	return nil
}
