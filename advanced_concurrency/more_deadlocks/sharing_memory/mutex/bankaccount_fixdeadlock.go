/* A primeira abordagem que podemos adotar é detectar deadlocks para que possamos
fazer algo a respeito. Por exemplo, após detectar que ocorreu um deadlock,
podemos ter um alerta que chama alguém que pode reiniciar o processo. Melhor ainda,
podemos ter uma lógica em nosso código que é notificada sempre que há um
impasse e executa uma operação de nova tentativa.

Go tem alguma detecção de deadlock integrada. O tempo de execução do Go verifica para ver
qual goroutine ele deve executar em seguida e, se descobrir que todos eles estão bloqueados
enquanto aguardam um recurso (como um mutex), ocorrerá um erro fatal. Infelizmente,
isso significa que ele só travará um deadlock se todas as goroutines estiverem bloqueadas.
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	accounts := []BankAccount2{
		*NewBankAccount2("Sam"),
		*NewBankAccount2("Paul"),
		*NewBankAccount2("Amy"),
		*NewBankAccount2("Mia"),
	}
	total := len(accounts)
	arb := NewArbitrator()
	for i := 0; i < 4; i++ {
		go func(tellerId int) {
			for i := 1; i < 1000; i++ {
				from, to := rand.Intn(total), rand.Intn(total)
				for from == to {
					to = rand.Intn(total)
				}
				accounts[from].Transfer2(&accounts[to], 10, tellerId, arb)
			}
			fmt.Println(tellerId, "COMPLETE")
		}(i)
	}
	time.Sleep(60 * time.Second)
}

type BankAccount2 struct {
	id      string
	balance int
	mutex   sync.Mutex
}

func NewBankAccount2(id string) *BankAccount2 {
	return &BankAccount2{
		id:      id,
		balance: 100,
		mutex:   sync.Mutex{},
	}
}

type Arbitrator struct {
	accountsInUse map[string]bool
	cond          *sync.Cond
}

func NewArbitrator() *Arbitrator {
	return &Arbitrator{
		accountsInUse: map[string]bool{},
		cond:          sync.NewCond(&sync.Mutex{}),
	}
}

func (a *Arbitrator) LockAccounts(ids ...string) {
	a.cond.L.Lock()
	for allAvailable := false; !allAvailable; {
		allAvailable = true
		for _, id := range ids {
			if a.accountsInUse[id] {
				allAvailable = false
				a.cond.Wait()
			}
		}
	}
	for _, id := range ids {
		a.accountsInUse[id] = true
	}
	a.cond.L.Unlock()
}

func (a *Arbitrator) UnlockAccounts(ids ...string) {
	a.cond.L.Lock()
	for _, id := range ids {
		a.accountsInUse[id] = false
	}
	a.cond.Broadcast()
	a.cond.L.Unlock()
}

func (src *BankAccount2) Transfer2(to *BankAccount2, amount int, tellerId int,
	arb *Arbitrator) {
	fmt.Printf("%d Locking %s and %s\n", tellerId, src.id, to.id)
	arb.LockAccounts(src.id, to.id)
	src.balance -= amount
	to.balance += amount
	arb.UnlockAccounts(src.id, to.id)
	fmt.Printf("%d Unlocked %s and %s\n", tellerId, src.id, to.id)
}
