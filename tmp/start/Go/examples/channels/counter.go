package main

import (
	"fmt"
	"time"
)

func sendData(ch chan int) {
	for i := 0; i < 100; i++ {
		ch <- i
	}
}

func getData(ch chan int) {
	for {
		if ch, ok := <-ch; ok {
			fmt.Println(ch)
		}
	}
}

func main() {
	ch := make(chan int)
	go sendData(ch)
	go getData(ch)
	time.Sleep(1 * time.Second)
}
