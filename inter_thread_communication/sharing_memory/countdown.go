/* A comunicação através da compartilhamento de memória é como tentar falar com um amigo,
mas em vez de trocar mensagens, usamos um quadro branco (ou um grande pedaço
de papel) e trocamos ideias, símbolos e abstrações

Na programação concorrente usando compartilhamento de memória, alocamos uma
parte da memória do processo – por exemplo, uma estrutura de dados compartilhada ou
uma variável – e temos diferentes goroutines trabalhando concorrentemente nesta
memória. Em nossa analogia, o quadro branco é a memória compartilhada usada
pelas diversas goroutines. */

package main

import (
	"fmt"
	"time"
)

/* Neste primeiro exemplo, criaremos
uma goroutine que compartilhará uma variável na memória com a goroutine main() (executando a
função main() ). A variável funcionará como um cronômetro de contagem regressiva. Uma goroutine
diminuirá o valor desta variável a cada segundo, e outra goroutine irá ler a variável com mais
frequência e exibi-la no console */

func main() {
	//cria uma variavel
	count := 5

	//cria uma goroutine passando a referencia dessa variavel
	go countdown(&count)

	for count > 0 {
		time.Sleep(500 * time.Millisecond)
		fmt.Println(count)
	}
}

func countdown(seconds *int) {
	for *seconds > 0 {
		time.Sleep(1 * time.Second)
		*seconds -= 1
	}
}

/* os dois for estao sendo executados ao mesmo tempo
na main de meio em meio segundo ele printa a variavel
e na countdown, de 1 em 1 segundo ele tira 1 da variavel */
