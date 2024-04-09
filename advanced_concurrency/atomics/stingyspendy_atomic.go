/* Os mutexes garantem que seções críticas do nosso código concorrente sejam
executadas por apenas uma goroutine por vez. Eles são usados para prevenir
race conditions. No entanto, os mutexes têm o efeito de transformar
partes da nossa programação concorrente em gargalos sequenciais.
Se estivermos apenas atualizando o valor de uma variável simples, como um
inteiro, podemos fazer uso de uma variável atômica para mantê-la consistente
entre goroutines sem precisar depender de mutexes que transformam nosso
código em um bloco sequencial.

Nos capítulos anteriores, vimos um exemplo com duas goroutines, chamadas Stingy
e Spendy, que compartilhavam uma variável inteira representando sua conta bancária.
O acesso à variável compartilhada foi protegido com um mutex. Este capítulo cobre
Cada vez que quiséssemos atualizar a variável, adquiriríamos o bloqueio mutex.
Assim que terminarmos a atualização, nós a lançaremos

Variáveis atômicas nos permitem realizar certas operações que são executadas sem interrupção.
Por exemplo, podemos adicionar o valor de uma variável compartilhada existente em uma única
operação atômica, o que garante que as operações de adição simultâneas não interfiram entre si.
Uma vez executada a operação, ela é totalmente aplicada ao valor da variável sem interrupção.
Podemos usar variáveis atômicas para substituir mutexes em determinados cenários.

Por exemplo, podemos facilmente alterar nosso programa Stingy e Spendy para usar essas operações
de variáveis atômicas. Em vez de usar mutexes, simplesmente chamaremos a operação atômica add()
em nossa variável de dinheiro compartilhado. Isso garante que as goroutines não produzam condições
de corrida que produzam resultados inconsistentes */

package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var money int32 = 100
	wg := sync.WaitGroup{}
	wg.Add(2)
	go stingy(&money, &wg)
	go spendy(&money, &wg)
	wg.Wait()
	fmt.Println("Money in account: ", atomic.LoadInt32(&money))
}

func stingy(money *int32, wg *sync.WaitGroup) {
	for i := 0; i < 1000000; i++ {
		atomic.AddInt32(money, 10)
	}
	fmt.Println("Stingy Done")
	wg.Done()
}

func spendy(money *int32, wg *sync.WaitGroup) {
	for i := 0; i < 1000000; i++ {
		atomic.AddInt32(money, -10)
	}
	fmt.Println("Spendy Done")
	wg.Done()
}

/* Por que não usamos operações atômicas para tudo, para eliminar o risco de
compartilhar uma variável e acidentalmente esquecendo de usar técnicas de
sincronização? Infelizmente, há uma penalidade de desempenho a ser paga sempre
que use essas variáveis atômicas. Atualizar uma variável de maneira normal é
um pouco mais rápido do que atualizar variáveis com operações atômicas.

Usar operações atômicas chega a ser três vezes mais lento do que usar o operador
normal. Esses resultados variarão em diferentes sistemas e arquiteturas, mas em
todos sistemas, há uma diferença substancial no desempenho.

Isso ocorre porque ao usar atômicos, estamos perdendo muitas otimizações de compilador
e sistema. Por exemplo, quando acessamos a mesma variável repetidamente, como fizemos,
o sistema mantém a variável no cache do processador, tornando o acesso à variável
mais rápido, mas pode liberar periodicamente a variável de volta para a memória
principal, especialmente se estiver ficando sem espaço em cache. Quando usando atômicos,
o sistema precisa garantir que qualquer outra execução em paralelo veja a atualização
para o variável. Assim, sempre que operações atômicas são utilizadas, o sistema precisa
manter as variáveis armazenadas em cache consistentemente. Isso pode ser feito liberando
a memória principal e invalidando quaisquer outros caches. Ter que manter vários caches
consistentes acabam reduzindo o desempenho do nosso programa. */
