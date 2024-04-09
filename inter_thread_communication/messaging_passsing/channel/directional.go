/* Os canais do Go são bidirecionais por padrão. Isso significa que uma goroutine pode atuar tanto como
receptor quanto como remetente de mensagens. No entanto, podemos atribuir uma direção a um canal
para que a goroutine que utiliza o canal possa apenas enviar ou receber mensagens.

No receptor, quando declaramos o canal como sendo mensagens <-chan int, estamos dizendo
que o canal é somente de recepção. A declaração de mensagens chan<-int na função sender diz o
contrário – que o canal só pode ser usado para enviar mensagens.*/

package main

import (
	"fmt"
	"time"
)

func main() {
	msgChannel := make(chan int)
	go receiver3(msgChannel)
	go sender(msgChannel)
	time.Sleep(5 * time.Second)
}

// recebe um output channel
func receiver3(messages <-chan int) {
	for {
		msg := <-messages
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
	}
}

// recebe um input channel
func sender(messages chan<- int) {
	for i := 0; ; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		messages <- i
		time.Sleep(1 * time.Second)
	}
}

/* se tentar enviar com método receptor ou vice versa dá erro na compilação */
