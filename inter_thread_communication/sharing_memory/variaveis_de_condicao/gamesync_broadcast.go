/* Vimos apenas exemplos usando Signal() em vez de
Broadcast() até agora. Quando temos várias goroutines suspensas
em Wait() de uma variável de condição, Signal() irá ativar
arbitrariamente uma dessas goroutines. A chamada Broadcast() ,
por outro lado, ativará todas as goroutines que estão suspensas em Wait().

O código abaixo exemplifica um jogo em que os jogadores esperam
que todos entrem antes do jogo começar. Este é um cenário comum tanto em
jogos multijogador online quanto em consoles de jogos. Vamos imaginar que
nosso programa tenha uma goroutine que trata das interações com cada jogador,
nosso código suspende a execução de cada goroutine
até que todos os jogadores entrem no jogo. */

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	cond := sync.NewCond(&sync.Mutex{})
	playersInGame := 4
	for playerId := 0; playerId < 4; playerId++ {
		go playerHandler(cond, &playersInGame, playerId)

		//Para simular cada jogador se conectando ao jogo em um momento diferente
		time.Sleep(1 * time.Second)
	}
}

func playerHandler(cond *sync.Cond, playersRemaining *int, playerId int) {
	cond.L.Lock()
	fmt.Println(playerId, ": Connected!")

	*playersRemaining--
	if *playersRemaining == 0 {
		cond.Broadcast() // acorda todas as goroutines suspensas de uma vez
	}

	for *playersRemaining > 0 {
		fmt.Println(playerId, ": Waiting for more players")
		cond.Wait()
	}

	cond.L.Unlock()
	fmt.Println("All players connected. Ready player", playerId)
}
