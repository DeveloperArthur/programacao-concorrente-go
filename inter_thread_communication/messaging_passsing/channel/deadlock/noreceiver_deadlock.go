/* O que aconteceria se uma goroutine enviasse uma mensagem em um canal sem que houvesse outra
goroutine para ler essa mensagem? Os canais do Go são síncronos por padrão, o que significa que a
goroutine remetente será bloqueada até que haja uma goroutine receptora pronta para consumir a
mensagem.

Como nossa goroutine receptora () termina após 5 segundos, nenhuma outra goroutine está disponível para consumir
mensagens do canal. O tempo de execução do Go percebe isso e gera o erro fatal. */

package main

import (
	"fmt"
	"time"
)

func main() {
	msgChannel := make(chan string)
	go receiver(msgChannel)
	fmt.Println("Sending HELLO...")
	msgChannel <- "HELLO"
	fmt.Println("Sending THERE...")
	msgChannel <- "THERE"
	fmt.Println("Sending STOP...")
	msgChannel <- "STOP"
}

func receiver(messages chan string) {
	time.Sleep(5 * time.Second)
	fmt.Println("Receiver slept for 5 seconds")
}

/* A ideia principal aqui é que, por padrão, os canais do Go são síncronos. Um remetente bloqueará se não houver
uma goroutine consumindo sua mensagem, e um destinatário bloqueará da mesma forma se não houver uma
goroutine enviando uma mensagem. */
