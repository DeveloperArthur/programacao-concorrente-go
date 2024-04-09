/* Vimos como os mutexes permitem que apenas uma
goroutine tenha acesso a um recurso compartilhado,
enquanto um mutex read-write nos permitem especificar
múltiplas leituras simultâneas, mas escritas exclusivas.

Os semáforos nos fornecem um tipo diferente de controle de concorrencia,
pois podemos especificar o número de execuções concorrentes permitidas.
Os semáforos também podem ser usados como blocos de construção para o
desenvolvimento de ferramentas de concorrência mais complexas.

Os mutexes nos oferecem uma maneira de permitir que apenas uma execução
aconteça por vez. E se precisarmos permitir que um número variável de execuções
aconteça concorrentemente? Semáforos nos permitem especificar
quantos goroutines podem acessar nosso recurso e limitar a concorrencia
isso nos permite limitar a carga em um sistema.

Para efeito de comparação:
Um mutex garante que apenas uma goroutine tenha acesso exclusivo,
enquanto um semáforo garante que no máximo N goroutines tenham acesso.
Na verdade, um mutex oferece a mesma funcionalidade de um semáforo onde
N tem o valor 1. Um semáforo de contagem nos permite a flexibilidade de
escolher qualquer valor de N.
(Um semáforo com apenas uma permissão é às vezes chamado de semáforo binário) */

package main

import "sync"

type Semaphore struct {
	permits int
	cond    *sync.Cond
}

// constructor
func NewSemaphore(n int) *Semaphore {
	return &Semaphore{
		permits: n,
		cond:    sync.NewCond(&sync.Mutex{}),
	}
}

func (rw *Semaphore) Acquire() {
	rw.cond.L.Lock()

	// se não tiver licença, suspende e fica em espera
	for rw.permits <= 0 {
		rw.cond.Wait()
	}

	rw.permits--
	rw.cond.L.Unlock()
}

func (rw *Semaphore) Release() {
	rw.cond.L.Lock()

	rw.permits++     // aumenta o número de licenças disponíveis em 1
	rw.cond.Signal() // sinaliza a gorotuine suspensa que mais uma licença está disponível

	rw.cond.L.Unlock()
}
