package main

import (
	"fmt"
	"time"
)

func main() {
	semaphore := NewSemaphore(6)
	for i := 0; i < 10; i++ {
		go foo(semaphore)
	}
	time.Sleep(5 * time.Second)
}

func foo(semaphore *Semaphore) {
	semaphore.Acquire()
	time.Sleep(2 * time.Second)
	fmt.Println(time.Now().String() + " WORKING")
	semaphore.Release()
}

/* dá pra ver certinho que 6 começam a trabalhar, e depois
que eles finalizam, os outros 4 começam */
