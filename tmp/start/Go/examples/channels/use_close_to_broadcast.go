package main

import (
	"fmt"
	"math/rand"
	"time"
)

func waiter(i int, block, done chan struct{}) {
	time.Sleep(time.Duration(rand.Intn(3000)) * time.Millisecond)
	fmt.Println(i, "waiting for close...")
	// Block and wait for close
	<-block
	fmt.Println(i, "done!")
	// Use it to pass the results
	done <- struct{}{}
}

func main() {
	block, done := make(chan struct{}), make(chan struct{})
	for i := 0; i < 4; i++ {
		go waiter(i, block, done)
	}
	time.Sleep(5 * time.Second)
	fmt.Println("CLOSE")
	close(block)
	for i := 0; i < 4; i++ {
		<-done
	}
}
