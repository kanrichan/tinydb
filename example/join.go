package main

import (
	tiny "github.com/Yiwen-Chan/tinydb"
)

type jointable[T any, S any] struct {
	A T
	B S
}

func Join[T any, S any](a tiny.Table[T], b tiny.Table[S], on func(T, S) bool, cond func(T, S) bool) ([]jointable[T, S], error) {
	da, err := a.Select(func(t T) bool { return true })
	if err != nil {
		return nil, err
	}
	db, err := b.Select(func(t S) bool { return true })
	if err != nil {
		return nil, err
	}
	var out = make([]jointable[T, S], 0)
	for i := range da {
		for j := range db {
			if on(da[i], db[j]) && cond(da[i], db[j]) {
				out = append(out, jointable[T, S]{da[i], db[j]})
			}
		}
	}
	return out, nil
}
