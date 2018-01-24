package main

//import "fmt"

/*

// Simple load balancer

func Run() {
	in, out := make(chan *Worker), make(chan *Worker)
	for i := 0; i < NumWorkers; i++; {
		go worker(in, out)
	}
	go sendLotsOfWork(in)
	receiveLotsOfResult(out)
}

func main() {
	fmt.Println("vim-go")
}

*/

type Request struct {
	fn func() int // The operation to perform
	c chan int // The channel to return the result
}

var workFn func() int

func requester(work chan<- Request) {
	c := make(chan int)
	for {
		// Kill some time (fake load)
		//Sleep(1)
		work <-Request{workFn, c}
		result := <-c
		_ = result
		//furtherProcess(result)
	}
}

type Worker struct {
	requests chan Request // work to do (buffered channel)
	pending int // count of pending tasks
	index int
}

func (w* Worker) work(done chan *Worker) {
	for {
		req := <-w.requests // Get requests from balancer
		req.c <-req.fn() // Call the fn and send the result
		done <-w // we've finished
	}
}

type Pool []*Worker


type Balancer struct {
	pool Pool
	done chan *Worker
}


// Balancer function
func (b *Balancer) balance(work chan Request) {
	for  {
		select {
		case req := <-work // received a request
			b.dispatch(req) // so send it to the worker
		case w := <-b.done // a worker has finished its job
			b.completed(w)
		}
	}
}


func (b *Balancer) dispatch(req Request) {
	// Grab the least loaded worker
	w := heap.Pop(&b.pool).(*Worker)
	// Send it to the task
	w.requests <- req
	// One more in its work queue
	w.pending++
	// Put it into its place on the heap
	heap.Push(&b.pool, w)
}


func (b *Balancer) completed(w *Worker) {
	// One fewer in th queue
	w.pending--
	// Remove it from heap
	heap.Remove(&b.pool, w.index)
	// Put it into its place
	heap.Push(&b.pool, w)
}
