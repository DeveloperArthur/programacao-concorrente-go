/* Para corrigir o problema de deadlocks iremos modificar
a função doWork() para que ela bloqueie o mutex antes de
chamar signal(). Isso garante que a goroutine main() esteja
em estado de espera */

package main

import (
	"fmt"
	"sync"
)

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	cond.L.Lock()
	for i := 0; i < 50000; i++ {
		go doWork(cond)
		fmt.Println("Waiting for child goroutine")
		cond.Wait()
		fmt.Println("Child goroutine finished")
	}
	cond.L.Unlock()
}

func doWork(cond *sync.Cond) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	cond.L.Lock()
	cond.Signal()
	cond.L.Unlock()
}

/* Agora funciona! Porque mesmo se a primeira goroutine
enviar um sinal, e a main ainda não estiver esperando, o main
chama cond.Wait() dentro de um loop for que executa 50.000 vezes.
Isso significa que, mesmo que a primeira goroutine envie o sinal
antes que o main chame cond.Wait(), as outras goroutines ainda não
terão a oportunidade de executar cond.Signal() até que o main chame
cond.Wait() novamente. */
