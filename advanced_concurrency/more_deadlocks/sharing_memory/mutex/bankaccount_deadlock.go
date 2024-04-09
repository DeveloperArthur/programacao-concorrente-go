/* Imagine que trabalhamos em um banco e temos a tarefa de implementar um
software que lê transações contábeis para transferir fundos de uma conta para
outra. Uma transação subtrai o saldo de uma fonte */

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	accounts := []BankAccount{
		*NewBankAccount("Sam"),
		*NewBankAccount("Paul"),
		*NewBankAccount("Amy"),
		*NewBankAccount("Mia"),
	}
	total := len(accounts)
	for i := 0; i < 4; i++ {
		go func(eId int) {
			for j := 1; j < 1000; j++ {
				from, to := rand.Intn(total), rand.Intn(total)
				for from == to {
					to = rand.Intn(total)
				}
				accounts[from].Transfer(&accounts[to], 10, eId)
			}
			fmt.Println(eId, "COMPLETE")
		}(i)
	}
	time.Sleep(60 * time.Second)
}

type BankAccount struct {
	id      string
	balance int
	mutex   sync.Mutex
}

func NewBankAccount(id string) *BankAccount {
	return &BankAccount{
		id:      id,
		balance: 100,
		mutex:   sync.Mutex{},
	}
}

func (src *BankAccount) Transfer(to *BankAccount, amount int, exId int) {
	fmt.Printf("%d Locking %s’s account\n", exId, src.id)
	src.mutex.Lock()
	fmt.Printf("%d Locking %s’s account\n", exId, to.id)
	to.mutex.Lock()
	src.balance -= amount
	to.balance += amount
	to.mutex.Unlock()
	src.mutex.Unlock()
	fmt.Printf("%d Unlocked %s and %s\n", exId, src.id, to.id)
}

/* Cada vez que executamos, obtemos resultados ligeiramente
diferentes, nem sempre resultando em um deadlock. Isto se deve à
natureza não determinística das execuções concorrente.

A partir de nossa saída, podemos observar que alguns goroutines estão bloqueando
algumas contas e tentando adquirir bloqueios em outras. O impasse em nosso
exemplo acontece entre as goroutines 0, 2 e 3.

solução: advanced_concurrency/more_deadlocks/sharing_memory/bankaccount_fixdeadlock.go
*/
