package main

import (
	"encoding/json"
	"os"
)

type StorageType map[string]map[int]map[string]interface{}

type Sto interface {
	Read() (StorageType, error)
	Write(StorageType) error
}

type Storage struct {
	Chan chan operate
}

type operate struct {
	Handle   func() (int, error)
	Callback chan callback
}

type callback struct {
	id  int
	err error
}

func NewStorage() (*Storage, error) {
	var sto Storage
	sto.Chan = make(chan operate)
	go func() {
		select {
		case op := <-sto.Chan:
			id, err := op.Handle()
			op.Callback <- callback{id: id, err: err}
		}
	}()
	return &sto, nil
}

type JSONStorage struct {
	Storage
	Handle *os.File
}

func NewJSONStorage(path string) (*JSONStorage, error) {
	sto, err := NewStorage()
	if err != nil {
		return nil, err
	}
	fi, err := os.OpenFile(path, os.O_CREATE, 0644)
	return &JSONStorage{Storage: *sto, Handle: fi}, err
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
	Storage
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
