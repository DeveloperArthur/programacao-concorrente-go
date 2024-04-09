package main

import (
	"fmt"
)

func main() {
	semaphore := NewSemaphore(0)
	for i := 0; i < 50000; i++ {
		go doWork(semaphore)
		fmt.Println("Waiting for child goroutine ")
		semaphore.Acquire()
		fmt.Println("Child goroutine finished")
	}
}

func doWork(semaphore *Semaphore) {
	fmt.Println("Work started")
	fmt.Println("Work finished")
	semaphore.Release()
}

/* Com semáforo inicializado com 0
Isso nos dá um sistema no qual chamar a função Release()
atua como nosso sinal de conclusão do trabalho. A função Acquire() atua então
como nosso Wait(). Neste sistema, não importa se chamamos Acquire()
antes ou depois da conclusão do trabalho, pois o semáforo mantém um
registro de quantas vezes Release () foi chamado usando a contagem de
permissões. Se chamarmos antes, a goroutine irá bloquear e aguardar o
sinal Release() . Se ligarmos depois, a goroutine retornará imediatamente,
pois há uma licença disponível.
Nesse programa estamos usando um semáforo para saber
quando uma goroutine é concluída.*/
