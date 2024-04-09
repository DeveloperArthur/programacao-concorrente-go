/* Outro cenário útil é o bloqueio apenas por um determinado período de tempo,
aguardando uma operação em um canal. Assim como nos dois exemplos
anteriores, queremos verificar se uma mensagem chegou em um canal, mas
queremos esperar alguns segundos para ver se uma mensagem chega, em vez de
desbloquear imediatamente e fazer outra coisa. Isto é útil em muitas situações em
que as operações do canal são sensíveis ao tempo. Considere, por exemplo, um
aplicativo de negociação financeira, onde, se não recebermos uma atualização
do preço das ações dentro de um intervalo de tempo, precisaremos gerar alertas.

Podemos implementar esse comportamento usando uma goroutine separada que
envia uma mensagem em um canal extra após um tempo limite especificado.
Podemos então usar esse canal extra em nossa instrução select , junto
com os outros canais. Isso nos dará o efeito de bloqueio na instrução select até
que qualquer um dos canais fique disponível ou o tempo limite ocorra */

package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	t, _ := strconv.Atoi(os.Args[1])
	messages := sendMsgAfter2(3 * time.Second)
	timeoutDuration := time.Duration(t) * time.Second
	fmt.Printf("Waiting for message for %d seconds...\n", t)

	// a partir do momento que você faz um select,
	//ele vai ficar ouvindo todos os canais selecionados...
	select {
	case msg := <-messages:
		fmt.Println("Message received:", msg)

	//time.After retorna um canal, entao conferimos se esse canal tem mensagem
	//uma mensagem nesse canal é retornado depois de um determinado tempo
	case tNow := <-time.After(timeoutDuration):
		fmt.Println("Timed out. Waited until:", tNow.Format("15:04:05"))
	}
}

func sendMsgAfter2(seconds time.Duration) <-chan string {
	messages := make(chan string)
	go func() {
		time.Sleep(seconds)
		messages <- "Hello"
	}()
	return messages
}
