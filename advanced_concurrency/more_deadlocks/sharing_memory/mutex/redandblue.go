/* Um impasse, em um programa simultâneo, ocorre quando as execuções são
bloqueadas indefinidamente, esperando umas pelas outras para liberar
recursos. Deadlocks são um efeito colateral indesejável de certos programas
simultâneos em que execuções simultâneas tentam adquirir acesso
exclusivo a vários recursos ao mesmo tempo.
Podemos ter um programa que roda sem problemas por um longo tempo e,
de repente, a execução é interrompida, sem motivo óbvio. */

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	lockA := sync.Mutex{}
	lockB := sync.Mutex{}
	go red(&lockA, &lockB)
	go blue(&lockA, &lockB)
	time.Sleep(20 * time.Second)
	fmt.Println("Done")
}

func red(lock1, lock2 *sync.Mutex) {
	for {
		fmt.Println("Red: solicitando lock1")
		lock1.Lock()
		fmt.Println("Red: peguei lock1, solicitando lock2")
		lock2.Lock()
		fmt.Println("Red: peguei os dois locks")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Red: Locks liberados")
	}
}
func blue(lock1, lock2 *sync.Mutex) {
	for {
		fmt.Println("Blue: solicitando lock2")
		lock2.Lock()
		fmt.Println("Blue: peguei lock2, solicitando lock1")
		lock1.Lock()
		fmt.Println("Blue: peguei os dois locks")
		lock1.Unlock()
		lock2.Unlock()
		fmt.Println("Blue: Locks liberados")
	}
}

/*
Depois que a goroutine blue() adquire o lock2, ela precisa solicitar o lock1.
A goroutine red() está mantendo o lock1 e precisa solicitar o lock2.
Cada goroutine mantém um lock e depois solicita o outro. Como o outro lock é
mantido por outra goroutine, o segundo lock nunca é adquirido. Isso cria a
situação de deadlocks em que as duas goroutines estarão sempre esperando que a
outra goroutine libere seu lock.

independente de quanto dure o loop, eles sempre vao entrar nesse looping infinito:
console:
Blue: peguei lock2, solicitando lock1
Red: peguei lock1, solicitando lock2

a goroutine blue adquire o lock2 e solicita o lock1
a goroutine red adquire o lock1 e solicita o lock2

as vezes o blue demora pra conseguir o lock2
pq o red fica em looping pegando os dois locks e os liberando
mas quando blue pega o lock2, ja era!
*/
