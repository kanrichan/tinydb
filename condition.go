package tinydb

func (req *Request) All() *Request {
	req.condition = func() func(doc Document) bool {
		return func(doc Document) bool {
			return true
		}
	}
	return req
}

func (req *Request) WhereEqual(key string, value float64) *Request {
	req.condition = func() func(doc Document) bool {
		return func(doc Document) bool {
			if _, ok := doc[key]; !ok {
				return false
			}
			if _, ok := doc[key].(float64); !ok {
				return false
			}
			return doc[key].(float64) == value
		}
	}
	return req
}

func (req *Request) WhereLess(key string, value float64) *Request {
	req.condition = func() func(doc Document) bool {
		return func(doc Document) bool {
			if _, ok := doc[key]; !ok {
				return false
			}
			if _, ok := doc[key].(float64); !ok {
				return false
			}
			return doc[key].(float64) < value
		}
	}
	return req
}

func (req *Request) WhereMore(key string, value float64) *Request {
	req.condition = func() func(doc Document) bool {
		return func(doc Document) bool {
			if _, ok := doc[key]; !ok {
				return false
			}
			if _, ok := doc[key].(float64); !ok {
				return false
			}
			return doc[key].(float64) > value
		}
	}
	return req
}
