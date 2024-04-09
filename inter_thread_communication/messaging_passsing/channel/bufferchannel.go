/* Embora os canais sejam síncronos, podemos configurá-los para que armazenem um número de mensagens antes
de serem bloqueados. Quando usamos um canal em buffer, a goroutine remetente não será
bloqueada enquanto houver espaço disponível no buffer.

Quando criamos um canal, podemos especificar sua capacidade de buffer. Então, sempre que uma
goroutine remetente escreve uma mensagem sem que nenhum receptor a consuma, o canal armazenará
a mensagem. Isso significa que enquanto houver espaço no buffer, nosso
remetente não bloqueará e não teremos que esperar que um destinatário leia a mensagem.

O canal continuará armazenando mensagens enquanto houver capacidade no buffer. Assim que o buffer estiver
cheio, o remetente irá bloquear novamente. Esse acúmulo de buffer de
mensagens também pode ocorrer se o receptor for lento e não consumir as mensagens com rapidez suficiente
para acompanhar o remetente.

Assim que uma goroutine receptora estiver disponível para consumir as mensagens, as mensagens serão enviadas
ao destinatário na mesma ordem em que foram enviadas. Isso acontece mesmo se a goroutine do remetente não
enviar mais mensagens novas. Enquanto houver mensagens no buffer, uma goroutine receptora não será bloqueada.

Assim que a goroutine receptora consumir todas as mensagens e o buffer estiver vazio, a goroutine receptora
será bloqueada novamente. Quando o buffer está vazio, um receptor será bloqueado se não tivermos um
remetente ou se o remetente estiver produzindo mensagens em uma taxa mais lenta do que o receptor
pode lê-las.

Vamos agora tentar isso na prática. O código abaixo mostra um receptor de mensagens lento que
consome mensagens do canal inteiro a uma taxa de uma por segundo. Usamos time.Sleep() para
desacelerar a goroutine. Assim que a goroutine receiver() recebe um valor -1 , ela para de receber
mensagens e chama Done() em um grupo de espera. */

package main

import (
	"fmt"
	"sync"
	"time"
)

// cria um canal em buffer e alimenta mensagens no canal em
// uma taxa mais rápida do que nosso leitor pode consumi-las
func main() {
	msgChannel := make(chan int, 3) // capacity buffer de 3 mensagens
	wGroup := sync.WaitGroup{}
	wGroup.Add(1)
	go receiver2(msgChannel, &wGroup)

	for i := 1; i <= 6; i++ {
		size := len(msgChannel) // retorna o capacity atual do buffer
		fmt.Printf("%s Sending: %d. Buffer Size: %d\n",
			time.Now().Format("15:04:05"), i, size)
		msgChannel <- i // envia i
	}

	msgChannel <- -1 // envia -1
	wGroup.Wait()    // só pra main esperar execucao da outra goroutine, inves de usarmos um sleep...
}

func receiver2(messages chan int, wGroup *sync.WaitGroup) {
	msg := 0
	for msg != -1 {
		time.Sleep(1 * time.Second)
		msg = <-messages
		fmt.Println("Received:", msg)
	}
	wGroup.Done()
}

/* Dá pra ver direitinho que main envia 3 mensagens e bloqueia
e main vai mandando de 1 em 1 conforme o receiver vai processando
de 1 em 1, sempre respeitando capacity do buffer, veja o output:
16:47:42 Sending: 1. Buffer Size: 0
16:47:42 Sending: 2. Buffer Size: 1
16:47:42 Sending: 3. Buffer Size: 2
16:47:42 Sending: 4. Buffer Size: 3
Received: 1
16:47:43 Sending: 5. Buffer Size: 3
Received: 2
16:47:44 Sending: 6. Buffer Size: 3
Received: 3
Received: 4
Received: 5
Received: 6
Received: -1

O receptor consumirá uma mensagem a cada segundo, liberando um espaço no buffer que o
remetente preencherá rapidamente.
*/
