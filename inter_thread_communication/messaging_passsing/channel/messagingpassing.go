/* Até agora, falamos sobre fazer com que nossos goroutines resolvam problemas
compartilhando memória e usando controles de sincronização para evitar que eles
se atropelem. A passagem de mensagens é outra maneira de habilitar a comunicação
entre threads, é uma forma diferente de comunicação das execuções.
Na passagem de mensagens, threads de execuções passam cópias de mensagens entre si
sempre que precisam se comunicar. Como essas execuções não compartilham memória,
eliminamos os riscos de muitos tipos de race conditions.

Go utiliza o CSP (documentado no README) e nos fornece o
conceito de canal, que permite que goroutines se conectem, sincronizem e
compartilhem mensagens entre si

Quando distribuímos aplicativos executados em várias máquinas, a passagem de mensagens é a principal forma de
comunicação. Como os aplicativos são executados em máquinas separadas e não compartilham memória,
eles compartilham informações enviando mensagens por meio de protocolos comuns, como HTTP.*/

package main

import "fmt"

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
	msg := ""
	for msg != "STOP" {
		msg = <-messages
		fmt.Println("Received:", msg)
	}
}
