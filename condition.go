package main

func (req *Request) All() *Request {
	req.condition = func(doc Document) bool {
		return true
	}
	return req
}

type Compare func(interface{}, interface{}) bool

func (req *Request) WhereFunc(fn func(doc Document) bool) *Request {
	req.condition = fn
	return req
}

func (req *Request) Where(compare Compare, key string, value interface{}) *Request {
	return req.WhereFunc(func(doc Document) bool {
		if _, ok := doc[key]; !ok {
			return false
		}
		return compare(doc[key], value)
	})
}

func Equal(a interface{}, b interface{}) bool {
	return a == b
}

func Less(a interface{}, b interface{}) bool {
	if _, ok := a.(float64); !ok {
		return false
	}
	if _, ok := b.(int); !ok {
		return false
	}
	return int(a.(float64)) < b.(int)
}

func More(a interface{}, b interface{}) bool {
	if _, ok := a.(float64); !ok {
		return false
	}
	if _, ok := b.(int); !ok {
		return false
	}
	return int(a.(float64)) > b.(int)
}
