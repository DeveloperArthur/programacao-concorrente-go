/* Altere o código inter_thread_communication/sharing_memory/variaveis_de_condicao/gamesync_broadcast.go
para que, ainda usando variáveis de condição, os jogadores esperem por um número fixo de
segundos. Se todos os jogadores não tiverem entrado nesse período, os
goroutines devem parar de esperar e deixar o jogo começar sem todos os jogadores.
Dica: tente usar outra goroutine com um temporizador de expiração */

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	cancel := false
	go timeout(cond, &cancel)
	playersInGame := 5 // mudou pra 5, mas o for vai só até 4
	for playerId := 0; playerId < 4; playerId++ {
		go playerHandler(cond, &playersInGame, playerId, &cancel)
		time.Sleep(1 * time.Second)
	}
	time.Sleep(60 * time.Second)
}

func playerHandler(cond *sync.Cond, playersRemaining *int, playerId int, cancel *bool) {
	cond.L.Lock()
	fmt.Println(playerId, ": Connected!")

	*playersRemaining--
	if *playersRemaining == 0 {
		cond.Broadcast() // nunca vai ser zero pq playersInGame é 5, e só 4 goroutines foram acionadas
	}

	for *playersRemaining > 0 && !*cancel {
		fmt.Println(playerId, ": Waiting for more players")
		cond.Wait() // os 4 vão travar aqui esperando, quando destravar vão cair no if, cancel vai estar true
	}

	cond.L.Unlock()

	if *cancel {
		fmt.Println(playerId, ": Game cancelled")
	} else {
		fmt.Println("All players connected. Ready player", playerId)
	}
}

func timeout(cond *sync.Cond, cancel *bool) {
	time.Sleep(10 * time.Second)
	cond.L.Lock()
	fmt.Println("timeout error")
	*cancel = true
	cond.Broadcast()
	cond.L.Unlock()
}
