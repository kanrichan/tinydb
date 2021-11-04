package main

import (
	"encoding/json"
	"os"
)

type StorageType map[string]map[int]map[string]interface{}

type Storage interface {
	Read() (StorageType, error)
	Write(StorageType) error
}

type BaseStorage struct {
}

func NewStorage() (*BaseStorage, error) {
	var sto BaseStorage
	return &sto, nil
}

type JSONStorage struct {
	BaseStorage
	Handle *os.File
}

func NewJSONStorage(path string) (*JSONStorage, error) {
	sto, err := NewStorage()
	if err != nil {
		return nil, err
	}
	fi, err := os.OpenFile(path, os.O_CREATE, 0644)
	return &JSONStorage{BaseStorage: *sto, Handle: fi}, err
}

func (sto *JSONStorage) Read() (StorageType, error) {
	var data StorageType
	dec := json.NewDecoder(sto.Handle)
	err := dec.Decode(&data)
	return data, err
}

func (sto *JSONStorage) Write(data StorageType) error {
	enc := json.NewEncoder(sto.Handle)
	return enc.Encode(data)
}

type MemoryStorage struct {
	BaseStorage
	Chan   chan int
	Memory StorageType
}

func NewMemoryStorage() (*MemoryStorage, error) {
	return &MemoryStorage{Memory: StorageType{}}, nil
}

func (sto *MemoryStorage) Read() (StorageType, error) {
	return sto.Memory, nil
}

func (sto *MemoryStorage) Write(data StorageType) error {
	sto.Memory = data
	return nil
}
