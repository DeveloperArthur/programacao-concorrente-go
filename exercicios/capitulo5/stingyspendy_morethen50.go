/* Na inter_thread_communication/sharing_memory/variaveis_de_condicao/stingyspendy_wait.go
a goroutine de Stingy está sinalizando sobre a variável
de condição toda vez que adicionamos dinheiro à conta bancária.
Você pode alterar a função para que ela sinalize apenas
quando houver US$ 50 ou mais na conta? */

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
		if *money >= 50 {
			cond.Signal()
		}
		cond.L.Unlock()
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int, cond *sync.Cond) {
	for i := 0; i < 200000; i++ {
		cond.L.Lock()

		for *money < 50 {
			cond.Wait()
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
