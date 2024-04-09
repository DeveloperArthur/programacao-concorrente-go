/* Considere o código a seguir. Você consegue encontrar a race condition neste
programa sem executar o race detector? Dica: Tente executar o programa
diversas vezes para ver se ele resulta em uma race condition */

package main

import (
	"fmt"
	"time"
)

func addNextNumber(nextNum *[101]int) {
	i := 0
	for nextNum[i] != 0 { // é diferente de zero?
		i++ //sim, entao soma 1 em i e volta pro loop
	}
	// nao é diferente de zero
	nextNum[i] = nextNum[i-1] + 1 // entao soma +1 nessa posicao baseado no valor anterior
	//isso vai fazer as posicoes sempre serem sequenciais 1, 2, 3, 4, 5... 100, 101
}

func main() {
	nextNum := [101]int{1} // primeira posicao é 1

	for i := 0; i < 100; i++ {
		go addNextNumber(&nextNum) // gerar 99 goroutines
	}
	for nextNum[100] == 0 { // isso vai rodar enquanto a posicao 100 do array for 0
		println("Waiting for goroutines to complete")
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Println(nextNum)
}

/* Resposta: Ele resulta em race condition sim, algumas vezes ele printa
o array certinho de 1 a 101 mas algumas vezes ele só fica printando
"Waiting for goroutines to complete" sem parar...

Eu entendi o caminho feliz: o segundo for do main vai rodar
enquanto a posicao 100 do array for 0, e vão ter 99 goroutines
rodando for do addNextNumber, todas as goroutines vão navegar
da primeira posicao até a posicao que é igual a zero, porque
o for checa: é diferente de zero?, sim, entao soma 1 em i e volta pro loop
se nao é diferente de zero entao soma +1 nessa posicao baseado no valor anterior
isso vai fazer as posicoes sempre serem sequenciais 1, 2, 3, 4, 5... 100, 101
todas as goroutines estao atualizando o slice ao mesmo tempo, navegando e atualizando

mas por alguma razão, as vezes as goroutines param no meio do caminho
isso ocasiona o loop infinito do "Waiting for goroutines to complete"
as vezes elas param faltando o 101 [... 93 94 95 96 97 98 99 100 0]
mas as vezes param em outros números aleatórios, como 76, 88...

eu não consigo encontrar um motivo claro do porque as vezes
as goroutines param no meio do caminho */
