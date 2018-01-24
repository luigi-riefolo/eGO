package main

func Query(conns []Conn, query string) Result {
	ch := make(chan Result, len(conns))
	for _, conn := range conns {
		go func(c Conn) {
			ch <- c.DoQuery(query)
		}(conn)
	}
	return <-ch
}
