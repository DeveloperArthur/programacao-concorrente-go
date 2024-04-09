/* Como podemos fazer com que uma goroutine responda a mensagens provenientes de
diferentes goroutines em vários canais? A instrução select de Go nos permite especificar múltiplas
operações de canal como casos separados e então executar um caso dependendo de qual canal
está pronto.

Vamos pensar em um cenário simples onde uma goroutine espera mensagens de canais
separados, mas não sabemos em qual canal a próxima mensagem será recebida.
A instrução select nos permite agrupar as operações de leitura em vários
canais juntos, bloqueando a goroutine até que uma mensagem
chegue em qualquer um dos canais

Assim que uma mensagem chega em qualquer um dos canais, a goroutine é
desbloqueada e um manipulador de código para aquele canal é executado */

package main

import (
	"fmt"
	"time"
)

func main() {
	messagesFromA := writeEvery("Tick", 1*time.Second)
	messagesFromB := writeEvery("Tock", 2*time.Second)

	for {
		//quando tiver mensagens disponiveis dos canais, serao printados
		select {
		case msg1 := <-messagesFromA:
			fmt.Println(msg1)
		case msg2 := <-messagesFromB:
			fmt.Println(msg2)
		}

		/*Ao usar select, se vários casos estiverem prontos, um caso será escolhido
		aleatoriamente. Seu código não deve depender da ordem em que os casos são
		especificados*/
	}
}

func writeEvery(msg string, seconds time.Duration) <-chan string {
	//cria um canal
	messages := make(chan string)

	//grava no canal infinitamente
	go func() {
		for {
			time.Sleep(seconds)
			messages <- msg
		}
	}()

	//e retorna o canal paralelamente
	return messages
}
