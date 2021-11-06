package main

import (
	"encoding/json"
	"os"
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
	Mutex  sync.Mutex
	Handle *os.File
}

func NewJSONStorage(path string) (*JSONStorage, error) {
	sto, err := NewBaseStorage()
	if err != nil {
		return nil, err
	}
	fi, err := os.OpenFile(path, os.O_CREATE, 0644)
	return &JSONStorage{BaseStorage: *sto, Handle: fi}, err
}

func (sto *JSONStorage) Read() (*StorageData, error) {
	var data = StorageData{}
	dec := json.NewDecoder(sto.Handle)
	err := dec.Decode(&data.data)
	return &data, err
}

func (sto *JSONStorage) Write(data *StorageData) error {
	sto.Mutex.Lock()
	defer sto.Mutex.Unlock()
	enc := json.NewEncoder(sto.Handle)
	return enc.Encode(data.data)
}

func (sto *JSONStorage) Close() error {
	return sto.Handle.Close()
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
