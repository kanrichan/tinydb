package main

import "errors"

func TinyDB(storage Storage, err error) (*database, error) {
	if err != nil {
		return nil, err
	}
	tabs, err := storage.Read()
	if err != nil {
		return nil, err
	}
	if _, ok := tabs.item["_default"]; !ok {
		tabs.Lock()
		defer tabs.Unlock()
		tabs.item = TinyTabs{}
		if err := storage.Write(tabs); err != nil {
			return nil, err
		}
	}
	return &database{"_default", storage}, nil
}

func (tiny *database) SetTable(name string) error {
	tabs, err := tiny.storage.Read()
	if err != nil {
		return err
	}
	if _, ok := tabs.item[name]; !ok {
		tabs.Lock()
		defer tabs.Unlock()
		tabs.item = TinyTabs{}
		if err := tiny.storage.Write(tabs); err != nil {
			return err
		}
	}
	tiny.table = name
	return nil
}

func (tiny *database) Insert(rec TinyRecs) error {
	tabs, err := tiny.storage.Read()
	if err != nil {
		return err
	}
	tabs.Lock()
	defer tabs.Unlock()
	if _, ok := tabs.item[tiny.table]; !ok {
		tabs.item[tiny.table] = TinyDocs{}
	}
	var max = 0
	for i := range tabs.item[tiny.table] {
		if max < i {
			max = i
		}
	}
	tabs.item[tiny.table][max+1] = rec
	return tiny.storage.Write(tabs)
}

func (tiny *database) Search(fn func(recs TinyRecs) bool) ([]TinyRecsMap, error) {
	tabs, err := tiny.storage.Read()
	if err != nil {
		return nil, err
	}
	if _, ok := tabs.item[tiny.table]; !ok {
		return nil, errors.New("No such table")
	}
	tabs.RLock()
	defer tabs.RUnlock()
	var recs = []TinyRecsMap{}
	for i := range tabs.item[tiny.table] {
		if !fn(tabs.item[tiny.table][i]) {
			continue
		}
		recs = append(recs, TinyRecsMap{tabs.Docs(tiny.table), i})
	}
	return recs, nil
}

func (tiny *database) Update(recs TinyRecs, fn func(recs TinyRecs) bool) error {
	tabs, err := tiny.storage.Read()
	if err != nil {
		return err
	}
	if _, ok := tabs.item[tiny.table]; !ok {
		return errors.New("No such table")
	}
	tabs.Lock()
	defer tabs.Unlock()
	for i := range tabs.item[tiny.table] {
		if !fn(tabs.item[tiny.table][i]) {
			continue
		}
		for k := range recs {
			tabs.item[tiny.table][i][k] = recs[k]
		}
	}
	return tiny.storage.Write(tabs)
}

func (tiny *database) Remove(fn func(recs TinyRecs) bool) error {
	tabs, err := tiny.storage.Read()
	if err != nil {
		return err
	}
	if _, ok := tabs.item[tiny.table]; !ok {
		return errors.New("No such table")
	}
	tabs.Lock()
	defer tabs.Unlock()
	for i := range tabs.item[tiny.table] {
		if !fn(tabs.item[tiny.table][i]) {
			continue
		}
		delete(tabs.item[tiny.table], i)
	}
	return tiny.storage.Write(tabs)
}

func Equal(key string, value int) func(recs TinyRecs) bool {
	return func(recs TinyRecs) bool {
		if _, ok := recs[key]; !ok {
			return false
		}
		if _, ok := recs[key].(float64); !ok {
			return false
		}
		return int(recs[key].(float64)) == value
	}
}
