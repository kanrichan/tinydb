package main

type TableType map[int]DocumentType

type DocumentType map[string]interface{}

type Request struct {
	table     string
	operation func() func(storage Storage) ([]DocumentType, error)
	condition func() func(doc DocumentType) bool
	response  chan *Response
}

type Response struct {
	err  error
	docs []DocumentType
}

func Insert(doc ...DocumentType) *Request {
	req := &Request{}
	req.operation = func() func(storage Storage) ([]DocumentType, error) {
		return func(storage Storage) ([]DocumentType, error) {
			data, err := storage.Read()
			if err != nil {
				return nil, err
			}
			if _, ok := data[req.table]; !ok {
				data[req.table] = TableType{}
			}
			var id int
			for i := range data[req.table] {
				if id < i {
					id = i
				}
			}
			for i := range doc {
				data[req.table][id+i+1] = doc[i]
			}
			return doc, storage.Write(data)
		}
	}
	return req
}

func Remove() *Request {
	req := &Request{}
	req.operation = func() func(storage Storage) ([]DocumentType, error) {
		return func(storage Storage) ([]DocumentType, error) {
			data, err := storage.Read()
			if err != nil {
				return nil, err
			}
			if _, ok := data[req.table]; !ok {
				data[req.table] = TableType{}
			}
			var docs []DocumentType
			for id, doc := range data[req.table] {
				if req.condition()(doc) {
					docs = append(docs, data[req.table][id])
					delete(data[req.table], id)
				}
			}
			return docs, storage.Write(data)
		}
	}
	return req
}

func Search() *Request {
	req := &Request{}
	req.operation = func() func(storage Storage) ([]DocumentType, error) {
		return func(storage Storage) ([]DocumentType, error) {
			data, err := storage.Read()
			if err != nil {
				return nil, err
			}
			if _, ok := data[req.table]; !ok {
				data[req.table] = TableType{}
			}
			var docs []DocumentType
			for id, doc := range data[req.table] {
				if req.condition()(doc) {
					docs = append(docs, data[req.table][id])
				}
			}
			return docs, nil
		}
	}
	return req
}

func Update(kv map[string]interface{}) *Request {
	req := &Request{}
	req.operation = func() func(storage Storage) ([]DocumentType, error) {
		return func(storage Storage) ([]DocumentType, error) {
			data, err := storage.Read()
			if err != nil {
				return nil, err
			}
			if _, ok := data[req.table]; !ok {
				data[req.table] = TableType{}
			}
			var docs []DocumentType
			for _, doc := range data[req.table] {
				if req.condition()(doc) {
					docs = append(docs, doc)
					for key := range kv {
						doc[key] = kv[key]
					}
				}
			}
			return docs, storage.Write(data)
		}
	}
	return req
}

func (req *Request) All() *Request {
	req.condition = func() func(doc DocumentType) bool {
		return func(doc DocumentType) bool {
			return true
		}
	}
	return req
}

func (req *Request) WhereEqual(key string, value float64) *Request {
	req.condition = func() func(doc DocumentType) bool {
		return func(doc DocumentType) bool {
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
	req.condition = func() func(doc DocumentType) bool {
		return func(doc DocumentType) bool {
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
	req.condition = func() func(doc DocumentType) bool {
		return func(doc DocumentType) bool {
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
