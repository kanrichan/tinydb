package tinydb

import (
	"encoding/json"
)

func CachingMiddleware(storage Storage, err error) (*MiddlewareCaching, error) {
	return &MiddlewareCaching{
		storage: storage,
		size:    1000,
	}, err
}

func (mv *MiddlewareCaching) Read() (TinyTabs, error) {
	if mv.cache == nil {
		return mv.storage.Read()
	}
	var tabs TinyTabs
	json.Unmarshal(mv.cache, &tabs)
	return tabs, nil
}

func (mv *MiddlewareCaching) Write(tabs TinyTabs) error {
	b, err := json.Marshal(tabs)
	if err != nil {
		return err
	}
	mv.cache = b
	mv.count++
	if mv.count >= mv.size {
		tab, err := mv.Read()
		if err != nil {
			return err
		}
		return mv.storage.Write(tab)
	}
	return nil
}

func (mv *MiddlewareCaching) Close() error {
	tab, err := mv.Read()
	if err != nil {
		return err
	}
	err = mv.storage.Write(tab)
	if err != nil {
		return err
	}
	return mv.storage.Close()
}
