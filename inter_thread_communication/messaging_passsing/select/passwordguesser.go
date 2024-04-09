/* Um cenário útil é usar o caso de seleção padrão para cálculos e então usar um
canal para sinalizar quando precisamos parar. Para ilustrar esse
conceito, suponha que temos um aplicativo de exemplo que descobrirá uma senha
esquecida por meio de força bruta. Para simplificar, digamos que temos um arquivo protegido
por senha que lembramos que tem uma senha de seis caracteres ou menos, usando apenas
as letras minúsculas de a a z e espaços.

O número de sequências possíveis de "a" a "zzzzzz", incluindo espaços, é 27 6 – 1
(387.420.488). A função na listagem a seguir nos dá uma maneira para converter os
inteiros de 1 a 387.420.488 em uma string. Por exemplo, chamar toBase27(1) nos dá “a”,
chamá-lo com 2 nos dá “b”, 28 nos dá “aa” e assim por diante. */

package main

import (
	"fmt"
	"time"
)

const (
	passwordToGuess = "go far"
	alphabet        = " abcdefghijklmnopqrstuvwxyz"
)

// Algorithm converts a decimal integer into a string of base 27 using the alphabet constant
func toBase27(n int) string {
	result := ""
	for n > 0 {
		result = string(alphabet[n%27]) + result
		n /= 27
	}
	return result
}

/* Se tivéssemos que usar uma abordagem de força bruta em um programa sequencial, apenas
criaríamos um loop enumerando todas as strings de "a" a "zzzzzz" e, sempre, verificaríamos
se correspondia à variável passwordToGuess. Num cenário da vida real, não
teríamos o valor da senha; em vez disso, tentaríamos obter acesso ao nosso recurso (como um
arquivo) usando cada enumeração de string como senha.

Para encontrar nossa senha mais rapidamente, podemos dividir o intervalo de nossas
suposições entre várias goroutines. Por exemplo, a goroutine A tentaria suposições a partir de
enumerações de strings de 1 a 10 milhões, a goroutine B tentaria suposições a partir de 10 milhões
para 20 milhões, e assim por diante. Dessa forma, podemos ter muitas
goroutines, cada uma trabalhando em uma parte separada do nosso espaço de problemas.*/

func guessPassword(from int, upto int, stop chan int, result chan string) {
	for guessN := from; guessN < upto; guessN += 1 {
		select {
		/* Dentro do loop sempre que tentam ler algo do canal stop
		nunca tem mensagem, porque em nenhum lugar gravamos mensagem
		nesse canal, mas depois que fechamos o canal, todas as tentativas
		de leitura de mensagens retornam 0, por isso todas as goroutines
		executam o case <-stop, depois de fechado ele começa a retornar 0*/
		case <-stop:
			fmt.Printf("Stopped at %d [%d,%d)\n", guessN, from, upto)
			return

		default:
			if toBase27(guessN) == passwordToGuess {
				result <- toBase27(guessN)
				close(stop)
				return
			}
		}
	}
	fmt.Printf("Not found between [%d,%d)\n", from, upto)
}

/* Para evitar cálculos desnecessários, queremos interromper a execução de cada goroutine
quando qualquer goroutine fizer uma estimativa correta. Para conseguir isso, podemos usar
um canal para notificar todas as outras goroutines quando uma execução descobrir a
senha. Depois que uma goroutine encontra a senha correspondente, ela fecha um canal
comum. Isso tem o efeito de interromper todas as goroutines participantes e
interromper o processamento.

Agora podemos criar diversas goroutines. Cada goroutine tentará encontrar
a senha correta dentro de um determinado intervalo.
A função main() cria os canais necessários e inicia todas as goroutines
com seus intervalos de entrada em passos de 10 milhões. */

func main() {
	finished := make(chan int)
	passwordFound := make(chan string)

	//Creates a goroutine with input ranges [1, 10M), [10M, 20M), . . . [380M, 390M)
	for i := 1; i < 387_420_488; i += 10_000_000 {
		go guessPassword(i, i+10_000_000, finished, passwordFound)
	}

	fmt.Println("password found:", <-passwordFound)
	close(passwordFound)
	time.Sleep(5 * time.Second)
}
