package main

type connect struct {
	table string
	conn  chan *Request
	stop  chan bool
}

func TinyDB(storage Storage) *connect {
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
	return &connect{"_default", conn, stop}
}

func (conn *connect) SetTable(name string) *connect {
	conn.table = name
	return conn
}

func (conn *connect) Exec(req *Request) *Response {
	response := make(chan *Response)
	req.table = conn.table
	req.response = response
	conn.conn <- req
	resp := <-req.response
	close(response)
	return resp
}

func (conn *connect) Close() {
	conn.stop <- true
	close(conn.conn)
	close(conn.stop)
}
