/* Podemos proteger seções críticas do nosso código com mutexes para que apenas uma goroutine por vez
acesse um recurso compartilhado. Desta forma, eliminamos as race conditions. Variações em
mutexes, às vezes chamadas de bloqueios, são usadas em todas as linguagens que suportam
programação concorrente

Se apenas uma goroutine acessar uma seção crítica por vez, estaremos protegidos das
race conditions. Afinal, as race conditions acontecem apenas quando há conflito
entre duas ou mais goroutines.

O código abaixo corrige, usando mutex, o programa que está em:
inter_thread_communication/sharing_memory/stingyspendy_race_condition.go */

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	money := 100
	mutex := sync.Mutex{}

	/* O uso do mutex garante que apenas uma goroutine por vez possa executar
	a seção crítica do código protegida pelo mutex. Portanto, mesmo que as duas
	goroutines executem métodos diferentes (stingy e spendy), o mutex garante que
	elas alternem de forma segura, garantindo exclusão mútua e evitando race condition */
	go stingy(&money, &mutex)
	go spendy(&money, &mutex)

	time.Sleep(2 * time.Second)

	/* Usamos um mutex quando lemos a variável money
	após o término das goroutines. Uma race condition aqui é muito improvável,
	já que fizemos sleep por 2 segundos para ter certeza de que as goroutines
	estão completas. No entanto, é sempre uma boa prática proteger os recursos
	partilhados, mesmo que tenha a certeza de que não haverá conflito */
	mutex.Lock()
	fmt.Println("Money in back account: ", money)
	mutex.Unlock()
}

func stingy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 1000000; i++ {
		mutex.Lock()
		*money += 10
		mutex.Unlock()
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, mutex *sync.Mutex) {
	for i := 0; i < 1000000; i++ {
		mutex.Lock()
		*money -= 10
		mutex.Unlock()
	}
	fmt.Println("Spendy Done")
}

/* IMPORTANTE: Devemos proteger todas as seções críticas, incluindo aquelas onde a
goroutine está apenas lendo os recursos compartilhados. As otimizações do
compilador podem reordenar as instruções, fazendo com que sejam executadas de
maneira diferente. O uso de mecanismos de sincronização adequados, como mutexes,
garante que estamos lendo a cópia mais recente dos recursos compartilhados */
