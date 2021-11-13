package tinydb

import (
	"math/rand"
	"regexp"
	"sort"
	"time"
)

func All() Selector {
	return func(docs TinyDocs) []int {
		var ls = []int{}
		for i := range docs {
			ls = append(ls, i)
		}
		return ls
	}
}

func Equal(key string, value interface{}) Selector {
	return func(docs TinyDocs) []int {
		var ls = []int{}
		for i := range docs {
			if _, ok := docs[i][key]; !ok {
				continue
			}
			if docs[i][key] == value {
				ls = append(ls, i)
			}
		}
		return ls
	}
}

func Like(key string, regex string) Selector {
	return func(docs TinyDocs) []int {
		var ls = []int{}
		for i := range docs {
			if _, ok := docs[i][key]; !ok {
				continue
			}
			if _, ok := docs[i][key].(string); !ok {
				continue
			}
			f, err := regexp.MatchString(regex, docs[i][key].(string))
			if err != nil || f {
				ls = append(ls, i)
			}
		}
		return ls
	}
}

func Less(key string, value float64) Selector {
	return func(docs TinyDocs) []int {
		var ls = []int{}
		for i := range docs {
			if _, ok := docs[i][key]; !ok {
				continue
			}
			if _, ok := docs[i][key].(float64); !ok {
				continue
			}
			if docs[i][key].(float64) < value {
				ls = append(ls, i)
			}
		}
		return ls
	}
}

func Greater(key string, value float64) Selector {
	return func(docs TinyDocs) []int {
		var ls = []int{}
		for i := range docs {
			if _, ok := docs[i][key]; !ok {
				continue
			}
			if _, ok := docs[i][key].(float64); !ok {
				continue
			}
			if docs[i][key].(float64) > value {
				ls = append(ls, i)
			}
		}
		return ls
	}
}

func NotLess(key string, value float64) Selector {
	return func(docs TinyDocs) []int {
		var ls = []int{}
		for i := range docs {
			if _, ok := docs[i][key]; !ok {
				continue
			}
			if _, ok := docs[i][key].(float64); !ok {
				continue
			}
			if docs[i][key].(float64) >= value {
				ls = append(ls, i)
			}
		}
		return ls
	}
}

func NotGreater(key string, value float64) Selector {
	return func(docs TinyDocs) []int {
		var ls = []int{}
		for i := range docs {
			if _, ok := docs[i][key]; !ok {
				continue
			}
			if _, ok := docs[i][key].(float64); !ok {
				continue
			}
			if docs[i][key].(float64) <= value {
				ls = append(ls, i)
			}
		}
		return ls
	}
}

func (s Selector) Limit(times int) Selector {
	return func(docs TinyDocs) []int {
		var ls = s(docs)
		if len(ls) >= times {
			return ls[:times]
		}
		return ls
	}
}

func (s Selector) OrderBy(key string, asc bool) Selector {
	return func(docs TinyDocs) []int {
		var ls = s(docs)
		if len(ls) <= 1 {
			return ls
		}
		sort.Slice(
			ls,
			func(i, j int) bool {
				if _, ok := docs[ls[i]]; !ok {
					return false
				}
				if _, ok := docs[ls[j]]; !ok {
					return false
				}
				if _, ok := docs[ls[i]][key]; !ok {
					return false
				}
				if _, ok := docs[ls[j]][key]; !ok {
					return false
				}
				if _, ok := docs[ls[i]][key].(float64); !ok {
					return false
				}
				if _, ok := docs[ls[j]][key].(float64); !ok {
					return false
				}
				return asc != (docs[ls[i]][key].(float64) > docs[ls[j]][key].(float64))
			},
		)
		return ls
	}
}

func (s Selector) OrderByRand() Selector {
	return func(docs TinyDocs) []int {
		var ls = s(docs)
		if len(ls) <= 1 {
			return ls
		}
		rand.Seed(time.Now().Unix())
		index := rand.Perm(len(ls))
		var new = make([]int, len(ls))
		for i := 0; i < len(ls); i++ {
			new[i] = ls[index[i]]
		}
		return new
	}
}
