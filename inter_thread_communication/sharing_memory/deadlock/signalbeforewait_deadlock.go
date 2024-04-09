/* O que acontece se uma goroutine chamar Signal() ou Broadcast() e não houver
nenhuma execução aguardando por ela? Ele será perdido ou armazenado para a
próxima goroutine chamar Wait()? Se não houver
nenhuma goroutine em estado de espera, a chamada Signal() ou Broadcast() será
perdida. Vejamos esse cenário usando variáveis de condição para resolver outro
problema: esperar que nossas goroutines concluam suas tarefas.

Até agora, temos usado time.Sleep() em nossa função main() para aguardar a
conclusão de nossas goroutines. Isso não é ótimo, já que estamos apenas estimando
quanto tempo as goroutines levarão. Se executarmos nosso código em um computador
mais lento, teremos que aumentar a quantidade de tempo que temos para dormir.

Em vez de usar sleep, podemos fazer com que nossa função main() espere por uma
variável de condição e então fazer com que a goroutine filha envie um sinal quando
estiver pronta */

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
	cond.Signal()
}

/* Ao executar nosso código, esse erro ocorre:
fatal error: all goroutines are asleep - deadlocks!

Esse problema ocorreu porque todas as goroutines de doWork já tinham
sinalizado mas main não tinha entrado em espera ainda, quando main
entra em espera, já é tarde... ela não recebe nenhum sinal, isso resulta deadlocks
ou seja, falta sincronização, main ficaria esperando eternamente
pois os sinais foram perdidos, precisamos fazer essa sincronização com mutex.

solução em: inter_thread_communication/sharing_memory/variaveis_de_condicao/signalbeforewait_wait.go */
