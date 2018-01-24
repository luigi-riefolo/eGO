package main

type struct Request{
	fn func() int // The operation to perform
	c chan int // The channel to return the result
}


func requester(work chan<- Request) {
	c := make(chan int)
	for {
		// Kill some time (fake load)
		Sleep(
		work <- Request{workFn, c}
		result := <- c
		furtherProcess(result)
	}
}

type Worker struct {
	request chan Request // work to do (buffered channel)
	pending int // count of pending tasks
	index int
}

func (w* Worker) work(done chan *Worker) {
	for {
		req := <- w.requests // Get requests from balancer
		req.c <- req.fn() //
		done <- w //
	}
}
