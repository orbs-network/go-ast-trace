package main

import "fmt"

func main() {
	ch := make(chan int)

	go channelWriteExample(ch)
	num := channelReadExample(ch)
	fmt.Println(num)
}

func channelReadExample(ch chan int) int {
	num := <- ch
	return num
}

func channelWriteExample(ch chan int) {
	ch <- 22
}