package main

func CachingMiddleware(storage Storage, err error) (*MiddlewareCaching, error) {
	return &MiddlewareCaching{
		storage: storage,
		size:    1000,
	}, err
}

func (mv *MiddlewareCaching) Read() (TinyTabsMap, error) {
	if mv.cache.item == nil {
		return mv.storage.Read()
	}
	return mv.cache, nil
}

func (mv *MiddlewareCaching) Write(data TinyTabsMap) error {
	mv.cache = data
	mv.count++
	if mv.count >= mv.size {
		return mv.storage.Write(mv.cache)
	}
	return nil
}

func (mv *MiddlewareCaching) Close() error {
	err := mv.storage.Write(mv.cache)
	if err != nil {
		return err
	}
	return mv.storage.Close()
}
