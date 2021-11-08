package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
)

type StorageType map[string]TableType

type Storage interface {
	Read() (StorageType, error)
	Write(StorageType) error
	Close() error
}

type storageJSON struct {
	mutex  sync.Mutex
	handle *os.File
}

func JSONStorage(file string) (*storageJSON, error) {
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
	return &storageJSON{handle: fi}, err
}

func (sto *storageJSON) Read() (StorageType, error) {
	var data = StorageType{}
	sto.handle.Seek(0, 0)
	dec := json.NewDecoder(sto.handle)
	err := dec.Decode(&data)
	if err != nil && err != io.EOF {
		return data, err
	}
	return data, nil
}

func (sto *storageJSON) Write(data StorageType) error {
	sto.mutex.Lock()
	defer sto.mutex.Unlock()
	sto.handle.Truncate(0)
	sto.handle.Seek(0, 0)
	enc := json.NewEncoder(sto.handle)
	enc.SetIndent("", "    ")
	if data == nil {
		return errors.New("Nothing needs to be written")
	}
	return enc.Encode(data)
}

func (sto *storageJSON) Close() error {
	return sto.handle.Close()
}

type storageMemory struct {
	Memory StorageType
}

func MemoryStorage() (*storageMemory, error) {
	return &storageMemory{Memory: StorageType{}}, nil
}

func (sto *storageMemory) Read() (StorageType, error) {
	return sto.Memory, nil
}

func (sto *storageMemory) Write(data StorageType) error {
	sto.Memory = data
	return nil
}

func (sto *storageMemory) Close() error {
	return nil
}
