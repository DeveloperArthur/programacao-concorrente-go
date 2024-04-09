/* No programa messagingpassing.go o receptor não emite a última mensagem STOP. Isso é
porque a goroutine main() termina antes que a goroutine receiver() tenha a chance de imprimir a
última mensagem. Você pode alterar a lógica, sem usar ferramentas extras de concorrencia e sem
usar a função sleep, para que a última mensagem seja impressa? */

package main

import (
	"fmt"
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

	<-msgChannel //espera uma mensagem chegar
	close(msgChannel)
}

func receiver(messages chan string) {
	msg := ""
	for msg != "STOP" {
		msg = <-messages
		fmt.Println("Received:", msg)
	}
	messages <- "" //manda uma mensagem
}

//ou seja, o main nao vai finalizar quando acabar a execucao
//ele vai esperar uma mensagem, o receiver vai enviar
//ai ele fecha o canal
