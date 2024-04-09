package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int)
	go worker(1, ch)
	go worker(2, ch)

	for i := 0; i < 10; i++ {
		ch <- i
	}

	time.Sleep(3 * time.Second)
}

func worker(workerId int, data chan int) {
	for msg := range data {
		time.Sleep(time.Second)
		fmt.Printf("Worker %d received %d\n", workerId, msg)
	}
}

/* Olha só que interessante, se eu tiver apenas 1 worker
o programa acaba em 10 segundos, de 1 em 1 segundo meu
único worker processa uma mensagem.
Com 2 workers o programa acaba em 5 segundos, porque
se o worker 1 pega uma mensagem, o worker 2 vai pegar
outra mensagem, a próxima! Então de 1 em 1 segundo
os 2 workers pegam 2 mensagens.
Com 3 workers o programa acaba em 4 segundos.
Com 5 workers o programa acaba em 2 segundos.
Com 10 workers o programa acaba em 1 segundo. */
