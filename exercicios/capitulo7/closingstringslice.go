/* No programa closing.go o receptor lê 0 quando o canal é fechado. Você pode tentar com
diferentes tipos de dados? O que acontece se o canal for do tipo string? E se for do tipo slice ?*/

package main

import (
	"fmt"
	"time"
)

func main() {
	msgChannel := make(chan []int)
	go recever(msgChannel)
	for i := 1; i <= 3; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		msgChannel <- make([]int, 10)
		time.Sleep(1 * time.Second)
	}
	close(msgChannel)
	time.Sleep(3 * time.Second)
}

func recever(messages <-chan []int) {
	for {
		msg, more := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg, more)
		time.Sleep(1 * time.Second)

		//verifica se o canal foi fechado
		if !more {
			return
		}
	}
}

/* string ele retorna uma string vazia
e slice ele retorna um slice vazio */
