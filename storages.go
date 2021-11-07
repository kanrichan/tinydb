package main

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
)

type StorageType map[string]TableType

type StorageData struct {
	mutex sync.Mutex
	data  StorageType
}

type Storage interface {
	Read() (*StorageData, error)
	Write(*StorageData) error
	Close() error
}

type BaseStorage struct {
}

func NewBaseStorage() (*BaseStorage, error) {
	return &BaseStorage{}, nil
}

type JSONStorage struct {
	BaseStorage
	mutex  sync.Mutex
	handle *os.File
}

func NewJSONStorage(file string) (*JSONStorage, error) {
	sto, err := NewBaseStorage()
	if err != nil {
		return nil, err
	}
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
	fi, err := os.OpenFile(file, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	return &JSONStorage{BaseStorage: *sto, handle: fi}, err
}

func (sto *JSONStorage) Read() (*StorageData, error) {
	var data = StorageType{}
	sto.mutex.Lock()
	defer sto.mutex.Unlock()
	sto.handle.Seek(0, 0)
	dec := json.NewDecoder(sto.handle)
	err := dec.Decode(&data)
	return &StorageData{data: data}, err
}

func (sto *JSONStorage) Write(data *StorageData) error {
	sto.mutex.Lock()
	defer sto.mutex.Unlock()
	sto.handle.Seek(0, 0)
	enc := json.NewEncoder(sto.handle)
	if data.data == nil {
		return errors.New("Nothing needs to be written")
	}
	return enc.Encode(data.data)
}

func (sto *JSONStorage) Close() error {
	return sto.handle.Close()
}

type MemoryStorage struct {
	BaseStorage
	Memory *StorageData
}

func NewMemoryStorage() (*MemoryStorage, error) {
	return &MemoryStorage{Memory: &StorageData{}}, nil
}

func (sto *MemoryStorage) Read() (*StorageData, error) {
	return sto.Memory, nil
}

func (sto *MemoryStorage) Write(data *StorageData) error {
	sto.Memory = data
	return nil
}

func (sto *MemoryStorage) Close() error {
	return nil
}
