package main

type BaseMiddleware struct {
	storage Storage
}

type CachingMiddleware struct {
	BaseMiddleware
	cache StorageType
	count int
	size  int
}

func NewCachingMiddleware(storage Storage, err error) (*CachingMiddleware, error) {
	return &CachingMiddleware{
		BaseMiddleware: BaseMiddleware{storage: storage},
		size:           1000,
	}, err
}

func (mv *CachingMiddleware) Read() (StorageType, error) {
	if mv.cache == nil {
		return mv.BaseMiddleware.storage.Read()
	}
	return mv.cache, nil
}

func (mv *CachingMiddleware) Write(data StorageType) error {
	mv.cache = data
	mv.count++
	if mv.count >= mv.size {
		return mv.storage.Write(mv.cache)
	}
	return nil
}

func (mv *CachingMiddleware) Close() error {
	err := mv.storage.Write(mv.cache)
	if err != nil {
		return err
	}
	return mv.storage.Close()
}
