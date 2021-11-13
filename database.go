package tinydb

import (
	"errors"
	"sync"
)

func TinyDB(storage Storage, err error) (*database, error) {
	if err != nil {
		return nil, err
	}
	var tiny = &database{sync.Mutex{}, "_default", storage}
	tiny.Lock()
	defer tiny.Unlock()
	tabs, err := storage.Read()
	if err != nil {
		return nil, err
	}
	if _, ok := tabs["_default"]; !ok {
		tabs["_default"] = TinyDocs{}
		if err := storage.Write(tabs); err != nil {
			return nil, err
		}
	}
	return tiny, nil
}

func (tiny *database) SetTable(name string) error {
	tiny.Lock()
	defer tiny.Unlock()
	tabs, err := tiny.storage.Read()
	if err != nil {
		return err
	}
	if _, ok := tabs[name]; !ok {
		tabs[name] = TinyDocs{}
		if err := tiny.storage.Write(tabs); err != nil {
			return err
		}
	}
	tiny.table = name
	return nil
}

func (tiny *database) DropTable(name string) error {
	tiny.Lock()
	defer tiny.Unlock()
	tabs, err := tiny.storage.Read()
	if err != nil {
		return err
	}
	delete(tabs, name)
	return tiny.storage.Write(tabs)
}

func (tiny *database) Insert(rec TinyRecs) error {
	tiny.Lock()
	defer tiny.Unlock()
	tabs, err := tiny.storage.Read()
	if err != nil {
		return err
	}
	if _, ok := tabs[tiny.table]; !ok {
		tabs[tiny.table] = TinyDocs{}
	}
	var max = 0
	for i := range tabs[tiny.table] {
		if max < i {
			max = i
		}
	}
	tabs[tiny.table][max+1] = rec
	return tiny.storage.Write(tabs)
}

func (tiny *database) Search(selector Selector) (*TinyRecsIter, error) {
	tiny.Lock()
	defer tiny.Unlock()
	tabs, err := tiny.storage.Read()
	if err != nil {
		return nil, err
	}
	if _, ok := tabs[tiny.table]; !ok {
		return nil, errors.New("No such table")
	}
	var recs = TinyRecsArr{}
	for _, i := range selector(tabs[tiny.table]) {
		recs = append(recs, tabs[tiny.table][i])
	}
	return recs.Iterator(), nil
}

func (tiny *database) Update(recs TinyRecs, selector Selector) error {
	tiny.Lock()
	defer tiny.Unlock()
	tabs, err := tiny.storage.Read()
	if err != nil {
		return err
	}
	if _, ok := tabs[tiny.table]; !ok {
		return errors.New("No such table")
	}
	for _, i := range selector(tabs[tiny.table]) {
		for k := range recs {
			tabs[tiny.table][i][k] = recs[k]
		}
	}
	return tiny.storage.Write(tabs)
}

func (tiny *database) Remove(selector Selector) error {
	tiny.Lock()
	defer tiny.Unlock()
	tabs, err := tiny.storage.Read()
	if err != nil {
		return err
	}
	if _, ok := tabs[tiny.table]; !ok {
		return errors.New("No such table")
	}
	for _, i := range selector(tabs[tiny.table]) {
		delete(tabs[tiny.table], i)
	}
	return tiny.storage.Write(tabs)
}
