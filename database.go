package main

func TinyDB(storage Storage) *database {
	var conn = make(chan *Request)
	var stop = make(chan bool)
	go func() {
		for {
			select {
			case req := <-conn:
				resp := Response{}
				resp.docs, resp.err = req.operation()(storage)
				req.response <- &resp
			case <-stop:
				storage.Close()
				return
			}
		}
	}()
	return &database{"_default", conn, stop}
}

func (conn *database) SetTable(name string) *database {
	conn.table = name
	return conn
}

func (conn *database) Exec(req *Request) *Response {
	response := make(chan *Response)
	req.table = conn.table
	req.response = response
	conn.conn <- req
	resp := <-req.response
	close(response)
	return resp
}

func (conn *database) Close() {
	conn.stop <- true
	close(conn.conn)
	close(conn.stop)
}
