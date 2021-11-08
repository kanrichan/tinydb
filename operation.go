package tinydb

func Insert(doc ...Document) *Request {
	req := &Request{}
	req.operation = func() func(storage Storage) ([]Document, error) {
		return func(storage Storage) ([]Document, error) {
			data, err := storage.Read()
			if err != nil {
				return nil, err
			}
			if _, ok := data[req.table]; !ok {
				data[req.table] = TableData{}
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
	req.operation = func() func(storage Storage) ([]Document, error) {
		return func(storage Storage) ([]Document, error) {
			data, err := storage.Read()
			if err != nil {
				return nil, err
			}
			if _, ok := data[req.table]; !ok {
				data[req.table] = TableData{}
			}
			var docs []Document
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
	req.operation = func() func(storage Storage) ([]Document, error) {
		return func(storage Storage) ([]Document, error) {
			data, err := storage.Read()
			if err != nil {
				return nil, err
			}
			if _, ok := data[req.table]; !ok {
				data[req.table] = TableData{}
			}
			var docs []Document
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
	req.operation = func() func(storage Storage) ([]Document, error) {
		return func(storage Storage) ([]Document, error) {
			data, err := storage.Read()
			if err != nil {
				return nil, err
			}
			if _, ok := data[req.table]; !ok {
				data[req.table] = TableData{}
			}
			var docs []Document
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
