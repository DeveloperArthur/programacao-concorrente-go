package main

import (
	"fmt"
	"runtime"
)

func sayHello() {
	fmt.Println("Started")
}

func main() {
	go sayHello()
	runtime.Gosched()
	fmt.Println("Finished")
}

/* Sem chamar o runtime.Gosched(), temos muito poucas chances de executar a
função sayHello(). A goroutine main() terminará antes que a goroutine que chama a
função sayHello() tenha tempo para ser executada na CPU. Como em Go saímos do
processo quando a goroutine main() termina, não veríamos o texto "Started" impresso.

Quando chamamos o agendador Go, ele pode pegar a outra goroutine e
começar a executá-la, ou pode continuar a execução da goroutine que chamou o
agendador

ao chamar runtime.Gosched() estamos apenas aumentando as chances
de que sayHello() seja executado.
Não há garantia de que isso acontecerá. */
