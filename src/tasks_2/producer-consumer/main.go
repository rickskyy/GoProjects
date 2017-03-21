package main

import (
	"fmt"
	"math/rand"

)

//Go implementation of Producer-Consumer problem"

type Storage struct {
	count int
	products []int
}

func Produce(storage *Storage, out chan <- int) {
	for i := 0; i < len(storage.products); i++ {

		//time.Sleep(1 * time.Millisecond)
		fmt.Println("Produced  ", storage.products[i])
		out <- storage.products[i]
	}
	close(out)
}

func Pack(in <- chan int, out chan <- int) {
	for i:= range in{
		fmt.Println("Pack", i)
		out <- i
	}
	close(out)
}

func Consume (in <- chan int) {

	sum := 0
	for i:= range in{
		sum += i
		fmt.Println("Consumed", i)
	}

	fmt.Println(sum)

}

func main() {
	const N = 200
	products := make([]int, N)
	storage := Storage{N, products}



	fromProduсer := make(chan int, 20)
	toConsumer := make(chan int, 100)

	for i := 0; i < N; i++ {
		storage.products[i] = rand.Intn(100)
	}

	go Produce(&storage, fromProduсer)
	go Pack(fromProduсer, toConsumer)
	Consume(toConsumer)
	Consume(toConsumer)
	Consume(toConsumer)
}

