/* Race conditions são o que acontece quando seu programa tenta fazer muitas coisas ao
mesmo tempo, e seu comportamento depende do tempo exato de eventos imprevisíveis independentes.
Como vimos na seção anterior, nosso programa de frequência de cartas acaba dando resultados
inesperados, mas às vezes o resultado é ainda mais dramático. Nosso código concorrente
pode funcionar perfeitamente por um longo período e, um dia, pode travar, resultando em corrupção de
dados mais grave. Isso pode acontecer porque as execuções concorrentes não têm sincronização
adequada e estão se sobrepondo */

/*Stingy e Spendy são duas goroutines separadas. Stingy trabalha duro e ganha o
dinheiro, mas nunca gasta um único dólar. Spendy é o oposto, gastar dinheiro sem
ganhar nada. Ambos os goroutines compartilham uma conta bancária comum.
Para demonstrar uma race condition, faremos com que Stingy e Spendy
ganhem e gastem 10 dólares de cada vez, 1 milhão de vezes. Como Spendy
está gastando exatamente a mesma quantia que Stingy está ganhando, devemos
terminar com a mesma quantia com que começamos se nossa programação
estiver correta */

package main

import (
	"fmt"
	"time"
)

func main() {
	money := 100
	go stingy(&money)
	go spendy(&money)
	time.Sleep(2 * time.Second)
	fmt.Println("Money in back account: ", money)
}

func stingy(money *int) {
	for i := 0; i < 1000000; i++ {
		*money += 10
	}
	fmt.Println("Stingy Done")
}

func spendy(money *int) {
	for i := 0; i < 1000000; i++ {
		*money -= 10
	}
	fmt.Println("Spendy Done")
}

/* Esperamos que o resultado seja 100 dólares. Afinal, estamos
apenas adicionando e subtraindo 10 à variável 1 milhão de vezes. Isso simula que
Stingy ganha 10 milhões e Spendy gasta a mesma quantia, deixando-nos com
o valor inicial de 100. No entanto, aqui está o resultado do programa:
Money in back account:  667580
Cada vez que executamos, há um valor inconsistente diferente
ficando até mesmo negativo: -9807230

Estamos tendo esse problema porque as operações *money += 10 e *money -= 10 não são
atômicas; após a compilação, eles se traduzem em mais de uma instrução. Uma interrupção
na execução pode ocorrer entre as instruções. Instruções diferentes de outra goroutine
podem interferir e causar race conditions. Quando esse excesso acontece, obtemos
resultados imprevisíveis

Uma seção crítica em nosso código é um conjunto de instruções que deve ser executado sem
interferência de outras execuções que afetem o estado usado naquela seção.
Quando esta interferência é permitida acontecer, podem surgir race conditions.

se setarmos apenas 1 CPU no inicio de main (runtime.GOMAXPROCS(1)) não
teremos race conditions, pois como vimos no README:
"o sistema operacional alternará a execução entre as goroutines de forma
que pareça haver paralelismo, mas na realidade as goroutines compartilham o mesmo recurso de CPU."

não vão rodar em paralelo, só serão alternadas, obviamente, esta
não é uma boa solução, principalmente porque estaríamos abrindo mão da
vantagem de ter múltiplos processadores, e também porque não há garantia de que
resolverá o problema completamente. Uma versão futura do Go pode
processar as goroutines com 1 CPU de maneira diferente e interromper nosso programa...

Independentemente do sistema de agendamento que utilizamos, devemos nos
proteger contra race conditions. Dessa forma, estamos protegidos de
problemas independente do ambiente em que o programa será executado.

Para evitar race conditions em nossa programação,
precisamos de uma boa sincronização e comunicação com o restante das goroutines para
garantir que elas não passem por cima umas das outras. Uma boa programação concorrente
envolve sincronizar efetivamente suas execuções concorrentes para eliminar condições de
corrida e, ao mesmo tempo, melhorar o desempenho e o rendimento

solução em:
- inter_thread_communication/sharing_memory/mutexes/stingyspendy_mutexes.go
- inter_thread_communication/sharing_memory/variaveis_de_condicao/stingyspendy_wait.go
*/
