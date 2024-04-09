package main

import (
	"fmt"
	"time"
)

func doWork(id int) {
	fmt.Printf("Work %d started at %s\n", id, time.Now().Format("15:04:05"))
	time.Sleep(1 * time.Second)
	fmt.Printf("Work %d finished at %s\n", id, time.Now().Format("15:04:05"))
}

//func main() {
//	for i := 0; i < 5; i++ {
//		doWork(i)
//	}
//}

/* Todo esse processo dura 5 segundos, podemos fazer de forma concorrente
adicionando a palavra chave go na frente da chamada da função doWork() */

func main() {
	for i := 0; i < 5; i++ {
		go doWork(i) //cria 1 goroutine toda vez que é chamado
		//a goroutine é executada de forma concorrente em uma
		//execução sepadada da função main()
	}
	time.Sleep(2 * time.Second)
}

/* goroutines não são threads, não threads como conhecemos,
são threads a nível de usuário ("user-level threads"),
é como ter diferentes threads rodando dentro da thread principal
em nível de kernel

Do ponto de vista do sistema operacional, um processo contendo
threads de nível de usuário parecerá ter apenas um thread de execução */
