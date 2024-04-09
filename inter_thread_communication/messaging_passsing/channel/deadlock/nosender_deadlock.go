/* Receptor bloqueado porque o remetente não está enviando nenhuma mensagem
também resulta em deadlocks */

package main

import (
	"fmt"
	"time"
)

func main() {
	msgChannel := make(chan string)
	go sender(msgChannel)
	fmt.Println("Reading message from channel...")
	msg := <-msgChannel
	fmt.Println("Received:", msg)
}

func sender(channel chan string) {
	time.Sleep(5 * time.Second)
	fmt.Println("Sender slept for 5 seconds")
}

/* A ideia principal aqui é que, por padrão, os canais do Go são síncronos. Um remetente bloqueará se não houver
uma goroutine consumindo sua mensagem, e um destinatário bloqueará da mesma forma se não houver uma
goroutine enviando uma mensagem. */
