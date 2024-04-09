/* Usando o canal quit para interromper a execução de uma goroutine
esse é um pattern muito utilizado, ter um channel quit */

package main

import "fmt"

func main() {
	numbers := make(chan int)
	quit := make(chan int)
	printNumbers(numbers, quit)
	next := 0
	for i := 1; ; i++ {
		next += i
		select {
		case numbers <- next:
		case <-quit:
			fmt.Println("Quitting number generation")
			return
		}
	}
}

func printNumbers(numbers <-chan int, quit chan int) {
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-numbers)
		}
		close(quit)
	}()
}
