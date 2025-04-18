package main

import "fmt"

func worker(ch chan<- int) {
	defer func() {
		fmt.Println("worker stopped")
	}()

	for i := 0; ; i++ {
		ch <- i
	}

}

func main() {
	ch := make(chan int)

	for i := 0; i < 3; i++ {
		go worker(ch)
	}

	for n := range ch {
		fmt.Println(n)
	}

	fmt.Println("exit")
}

// после получения n == 10 нужно что бы было выведено следующее
// worker stopped
// worker stopped
// worker stopped
// exit
