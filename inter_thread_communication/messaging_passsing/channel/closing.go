/* Go nos permite fechar um canal. Podemos fazer isso em código
chamando a função close(channel) . Depois de fecharmos um canal,
não devemos enviar mais mensagens para ele, pois isso gera erros.
Se tentarmos receber mensagens de um canal fechado, receberemos
mensagens contendo o valor padrão para o tipo de dados do canal.
Por exemplo, se nosso canal for do tipo inteiro, a leitura de um
canal fechado fará com que a operação de leitura retorne um valor 0. */

package main

import (
	"fmt"
	"time"
)

// envia algumas mensagens no canal, após
// as quais fecha o canal
func main() {
	msgChannel := make(chan int)
	go receiver4(msgChannel)
	for i := 1; i <= 3; i++ {
		fmt.Println(time.Now().Format("15:04:05"), "Sending:", i)
		msgChannel <- i
		time.Sleep(1 * time.Second)
	}
	close(msgChannel)
	time.Sleep(3 * time.Second)
}

// um receptor que consome mensagens continuamente mesmo depois
// de fecharmos o canal, com um loop que lê mensagens do canal
// e as envia para o console a cada segundo
func receiver4(messages <-chan int) {
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

// Essa é uma forma alternativa de consumir as mensagens sem validar se
// o canal está fechado, pois estamos continuamente iterando sob o canal,
// quando o canal fechar, o for simplesmente para
func newReceiver(messages <-chan int) {
	for msg := range messages {
		fmt.Println(time.Now().Format("15:04:05"), "Received:", msg)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("Receiver finished.")
}
