package main

type Middleware struct {
	storage Storage
}

type CachingMiddleware struct {
	Middleware
	cache StorageType
	count int
	size  int
}

func NewCachingMiddleware(storage Storage) *CachingMiddleware {
	return &CachingMiddleware{
		Middleware: Middleware{storage: storage},
	}
}

func (mi *CachingMiddleware) Read() (StorageType, error) {
	if mi.cache == nil {
		return mi.Middleware.storage.Read()
	}
	return mi.cache, nil
}

func (mi *CachingMiddleware) Write(data StorageType) error {
	mi.cache = data
	mi.count++
	if mi.count >= mi.size {
		return mi.storage.Write(mi.cache)
	}
	return nil
}
