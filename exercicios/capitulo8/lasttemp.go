/* No programa abaixo, temos duas goroutines. A função generateTemp() simula a leitura
e envio da temperatura em um canal a cada 200 ms. A função outputTemp()
simplesmente exibe uma mensagem encontrada em um canal a cada 2 segundos.
Você pode escrever uma função main() , usando uma instrução select , que leia
mensagens provenientes da goroutine generateTemp() e envie apenas
a temperatura mais recente para o canal outputTemp() ? Como a função
generateTemp() gera valores mais rapidamente do que a função outputTemp() ,
você precisará descartar alguns valores para que apenas a temperatura mais
atualizada seja exibida */

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	temps := generateTemp()
	display := make(chan int)
	outputTemp(display)
	t := <-temps
	for {
		select {
		case t = <-temps:
		case display <- t:
		}
	}
}

func generateTemp() chan int {
	output := make(chan int)
	go func() {
		temp := 50 //fahrenheit
		for {
			output <- temp
			temp += rand.Intn(3) - 1
			time.Sleep(200 * time.Millisecond)
		}
	}()
	return output
}

func outputTemp(input chan int) {
	go func() {
		for {
			fmt.Println("Current temp:", <-input)
			time.Sleep(2 * time.Second)
		}
	}()
}
