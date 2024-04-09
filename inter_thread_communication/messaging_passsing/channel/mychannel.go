package main

import (
	"container/list"
	"sync"
)

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

//-----------------------------------------------------------------------------------------------------------------

type Channel[M any] struct {
	capacitySema *Semaphore //Semáforo de capacidade para bloquear o remetente quando o buffer estiver cheio
	sizeSema     *Semaphore //Semáforo de tamanho de buffer para bloquear o receptor quando o buffer está vazio
	mutex        sync.Mutex //Mutex protegendo nossa estrutura de dados de lista compartilhada
	buffer       *list.List //Lista vinculada para ser usada como estrutura de dados de fila
}

func NewChannel[M any](capacity int) *Channel[M] {
	return &Channel[M]{
		capacitySema: NewSemaphore(capacity), //Cria um novo semáforo com o número de licenças igual à capacidade de entrada
		sizeSema:     NewSemaphore(0),        //Cria um novo semáforo com número de permissões igual a 0
		buffer:       list.New()}             //Cria uma nova lista vinculada vazia
}

func (c *Channel[M]) Send(message M) {
	c.capacitySema.Acquire() //Adquire uma licença do semáforo de capacidade

	//Adiciona uma mensagem à fila de buffer enquanto protege contra
	//race conditions usando um mutex
	c.mutex.Lock()
	c.buffer.PushBack(message)
	c.mutex.Unlock()

	//Libera uma permissão do semáforo de tamanho do buffer
	c.sizeSema.Release()
}

func (c *Channel[M]) Receive() M {
	c.capacitySema.Release() //Libera uma permissão do semáforo de capacidade

	c.sizeSema.Acquire() //Adquire uma licença do semáforo de tamanho de buffer

	//Remove uma mensagem do buffer enquanto protege contra
	//race condition usando o mutex
	c.mutex.Lock()
	v := c.buffer.Remove(c.buffer.Front()).(M)
	c.mutex.Unlock()

	//Retorna o valor da mensagem
	return v
}
