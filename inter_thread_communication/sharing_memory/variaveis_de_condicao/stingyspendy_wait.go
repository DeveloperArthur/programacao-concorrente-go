/* Este é o mesmo código que o stingyspendy_mutexes.go mas nesse caso
se a conta ficar negativa, o programa se encerra... E enquanto stingy poe
10 na conta, spendy tira 50, isso faz com que a conta fique negativa
muito rápido, não conseguimos concluir o programa... Precisamos encontrar
um jeito de parar a goroutine de spendy quando ela for deixar o saldo negativo,
e só voltar a trabalhar quando não for deixar o saldo negativo!
Para isso podemos usar variáveis de condição, elas funcionam junto com mutexes e
nos dão a capacidade de suspender a execução atual até termos um sinal de que
uma condição específica foi alterada...
No Go, usamos o pattern monitor sempre que usamos um mutex com uma
variável de condição */

package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

func main() {
	money := 100
	mutex := sync.Mutex{}
	variavelDeCondicao := sync.NewCond(&mutex)

	go stingy(&money, variavelDeCondicao)
	go spendy(&money, variavelDeCondicao)

	time.Sleep(2 * time.Second)

	mutex.Lock()
	fmt.Println("Money in back account: ", money)
	mutex.Unlock()
}

func stingy(money *int, cond *sync.Cond) {
	for i := 0; i < 1000000; i++ {
		cond.L.Lock()
		*money += 10
		cond.Signal() // enviando um sinal para spendy que valor foi adicionado
		cond.L.Unlock()
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, cond *sync.Cond) {
	for i := 0; i < 200000; i++ {
		cond.L.Lock()

		for *money < 50 {
			cond.Wait() // Goroutine é suspensa e coloca o mutex em espera

			//Após o sinal de stingy, a goroutine que estava suspensa volta daqui
			//se money > 50, sai do for e continua execução
		}

		*money -= 50
		if *money < 0 {
			fmt.Println("Money is negative!")
			os.Exit(1)
		}
		cond.L.Unlock()
	}
	fmt.Println("Spendy Done")
}

/* Se não houver fundos suficientes, a goroutine espera,
suspendendo a sua execução até que mais dinheiro
esteja disponível. Quando Stingy adiciona dinheiro,
ele envia um sinal para retomar qualquer execução que
esteja aguardando mais fundos.

exemplo:
começa com 100
spendy pega 50, fica 50
stingy poe 10, fica 60
spendy pega 50, fica 10
stingy poe 10, fica 20
spendy vê que money < 50 e dorme
stingy poe 10, fica 30
stingy poe 10, fica 40
stingy poe 10, fica 50
spendy acorda porque money >= 50
spendy pega 50, fica 0 */

/* Quando uma goroutine dá Unlock antes da finalização do método, ela simplesmente libera
o mutex para que outras goroutines possam adquiri-lo. A goroutine que deu Unlock continua
sua execução normalmente a partir desse ponto. O fato de liberar o mutex não interrompe a
execução da goroutine, ela continua até o final do método ou até encontrar outro ponto onde
possa ser bloqueada novamente pelo mutex, isso significa que teríamos  ambas as goroutines
executando simultaneamente. Isso é fundamental para a concorrência controlada por mutexes,
onde várias goroutines podem trabalhar em partes diferentes do código de forma concorrente,
mas com segurança, garantindo que apenas uma goroutine por vez acesse uma seção crítica. */
