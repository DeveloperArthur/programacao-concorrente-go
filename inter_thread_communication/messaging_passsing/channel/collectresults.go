/* Podemos executar funções concorrentemente em segundo plano e coletar seus resultados por meio de
canais assim que terminarem. Normalmente, na programação sequencial normal, chamamos uma função
e esperamos que ela retorne um resultado. Na programação concorrente, podemos chamar funções em
goroutines separadas e posteriormente obter seus valores de retorno em um canal de saída.

O programa a seguir mostra uma função que encontra os fatores de
um número de entrada. Por exemplo, se chamarmos findFactors(6)
ele retornará os valores [1 2 3 6] */

package main

import "fmt"

//func main() {
//	go fmt.Println(findFactors(3419110721))
//	fmt.Println(findFactors(4033836233))
//}

func findFactors(number int) []int {
	result := make([]int, 0)
	for i := 1; i <= number; i++ {
		if number%i == 0 {
			result = append(result, i)
		}
	}
	return result
}

/* Se chamarmos a função findFactors() duas vezes para
dois números diferentes sem "go", na programação sequencial,
teríamos duas chamadas, uma após a outra (main comentado).

Encontrar fatores de grandes números pode ser uma operação
demorada, por isso seria bom distribuir o trabalho em
vários núcleos de processamento.

Se tivermos vários núcleos disponíveis, executar a primeira
chamada findFactors() em paralelo com a segunda irá acelerar
nosso programa (main comentado)

Mas mesmo executando em paralelo, o problema é que se a
goroutine main finalizar primeiro, não veremos o resultado
da outra goroutine.

Poderíamos usar algo como uma variável compartilhada ou um
wait group, mas existe uma maneira mais fácil: usar canais */

func main() {
	resultCh := make(chan []int)

	//executando uma expressao de forma concorrente
	go func() {
		//gravando mensagem no canal
		resultCh <- findFactors(3419110721)
	}()

	fmt.Println(findFactors(4033836233))

	//lendo mensagem do canal
	fmt.Println(<-resultCh)
}
