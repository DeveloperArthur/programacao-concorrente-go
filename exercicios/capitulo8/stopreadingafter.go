/* No código abaixo temos uma função goroutine
generateNumbers() que gera números aleatórios. Você pode
escrever uma função main() usando uma instrução select que consome
continuamente do canal de saída, imprimindo a saída no console até
que 5 segundos tenham decorrido desde o início do programa?
Após 5 segundos, a função deverá parar de consumir do canal de
saída e o programa deverá ser encerrado. */

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	channel := generateNumbers()
	channelTime := time.After(5 * time.Second)
	for {
		select {
		case msg := <-channel:
			fmt.Println("Msg: ", msg)
		case <-channelTime:
			return
		}
	}
}

func generateNumbers() chan int {
	output := make(chan int)
	go func() {
		for {
			output <- rand.Intn(10)
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return output
}
